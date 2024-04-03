package services

import (
	"context"
	"time"

	galileo "github.com/Sharefunds/galileo-client"
	"github.com/Sharefunds/pool-transactions/pkg/client"
	"github.com/Sharefunds/pool-transactions/pkg/config"
	"github.com/Sharefunds/pool-transactions/pkg/db"
	"github.com/Sharefunds/pool-transactions/pkg/models"
	"github.com/Sharefunds/pool-transactions/pkg/utils"
	"github.com/antihax/optional"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Sharefunds/pool-transactions/pkg/pb"
)

type Server struct {
	H db.Handler
	pb.UnimplementedPoolTransactionsServiceServer
	GalileoClient        *galileo.APIClient
	C                    config.Config
	UserComplianceClient client.UserComplianceServiceClient
}

func (s *Server) CreatePool(ctx context.Context, req *pb.CreatePoolRequest) (*pb.CreatePoolResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid. - ")
	}

	userPrn, err := s.UserComplianceClient.GetUserPRN(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get PRN from User Compliance Service.")
	}

	opts := galileo.AccountsAndCardsApiPostAddaccountOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		ProdId:              optional.NewInt32(s.C.GalileoProductId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(userPrn.Prn),
		SharedBalance:       optional.NewInt32(0),
	}

	response, _, err := s.GalileoClient.AccountsAndCardsApi.PostAddaccount(ctx, &opts)

	if response.Status != "Success" {
		return nil, status.Errorf(codes.Internal, "Could not create account in Galileo.")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not create pool in Galileo.")
	}

	tx := s.H.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	PoolUUID := uuid.New()

	newPool := &models.Pool{
		PoolID:      PoolUUID,
		PRN:         response.ResponseData.PmtRefNo,
		PoolName:    req.Poolname,
		Description: req.Description,
		UserID:      userId,
		BalanceGoal: float64(req.BalanceGoal),
		Status:      "active",
	}

	if err := tx.Create(&newPool).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Could not create pool in database: %v", err)
	}

	participantUUID := uuid.New()

	newParticipant := &models.Participant{
		ParticipantID: participantUUID,
		PoolID:        newPool.PoolID,
		UserID:        userId,
		Role:          "owner",
		JoinedAt:      time.Now(),
	}

	if err := tx.Create(&newParticipant).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Could not add user as a participant in the pool: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Transaction commit error: %v", err)
	}

	return &pb.CreatePoolResponse{
		PoolId:      newPool.PoolID.String(),
		Poolname:    newPool.PoolName,
		Description: newPool.Description,
		BalanceGoal: int32(newPool.BalanceGoal),
	}, nil
}

func (s *Server) CreatePoolInvitation(ctx context.Context, req *pb.CreatePoolInvitationRequest) (*pb.CreatePoolInvitationResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	if req.PoolId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Pool ID and Invitee ID must be provided.")
	}

	var pool models.Pool
	if err := s.H.DB.Where("pool_id = ?", req.PoolId).First(&pool).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find pool with ID %s: %v", req.PoolId, err)
	}

	if pool.UserID != userId {
		return nil, status.Errorf(codes.PermissionDenied, "User is not the owner of the pool with ID %s", req.PoolId)
	}

	if pool.Status != "active" {
		return nil, status.Errorf(codes.PermissionDenied, "Pool with ID %s is not active", req.PoolId)
	}

	invitationUUID := uuid.New()

	invitation := &models.PoolInvitation{
		InvitationID: invitationUUID,
		PoolID:       req.PoolId,
		InviterID:    userId,
		InviteeID:    req.UserId,
		Status:       "pending",
	}

	result := s.H.DB.Where(models.PoolInvitation{
		PoolID:    req.PoolId,
		InviterID: userId,
		InviteeID: req.UserId,
	}).FirstOrCreate(invitation)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "Error creating invitation: %v", result.Error)
	}

	createdAt := timestamppb.New(invitation.CreatedAt)

	return &pb.CreatePoolInvitationResponse{
		InvitationId: invitation.InvitationID.String(),
		InviteeId:    invitation.InviteeID,
		PoolId:       invitation.PoolID,
		CreatedAt:    createdAt,
	}, nil
}

func (s *Server) CreatePoolParticipant(ctx context.Context, req *pb.CreatePoolParticipantRequest) (*pb.CreatePoolParticipantResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	if req.PoolId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Pool ID must be provided.")
	}

	var pool models.Pool
	if err := s.H.DB.Where("pool_id = ?", req.PoolId).First(&pool).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find pool with ID %s: %v", req.PoolId, err)
	}

	if pool.Status != "active" {
		return nil, status.Errorf(codes.PermissionDenied, "Pool with ID %s is not active", req.PoolId)
	}

	var poolInvitation models.PoolInvitation
	if err := s.H.DB.Where("pool_id = ? AND invitee_id = ?", req.PoolId, userId).First(&poolInvitation).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find invitation for pool with ID %s and user with ID %s: %v", req.PoolId, userId, err)
	}

	if poolInvitation.Status != "pending" {
		return nil, status.Errorf(codes.PermissionDenied, "Invitation for pool with ID %s and user with ID %s is not pending", req.PoolId, userId)
	}

	tx := s.H.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	participantUUID := uuid.New()

	newParticipant := &models.Participant{
		ParticipantID: participantUUID,
		PoolID:        uuid.MustParse(req.PoolId),
		UserID:        userId,
		Role:          "member",
		JoinedAt:      time.Now(),
	}

	if err := tx.Model(&poolInvitation).Update("status", "accepted").Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Could not update invitation status: %v", err)
	}

	if err := tx.Create(&newParticipant).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Could not add user as a participant in the pool: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Transaction commit error: %v", err)
	}

	return &pb.CreatePoolParticipantResponse{
		ParticipantId: newParticipant.ParticipantID.String(),
		PoolId:        newParticipant.PoolID.String(),
		Role:          newParticipant.Role,
		JoinedAt:      timestamppb.New(newParticipant.JoinedAt),
	}, nil
}

func (s *Server) DeletePool(ctx context.Context, req *pb.DeletePoolRequest) (*pb.DeletePoolResponse, error) {
	userID, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	var pool models.Pool
	if err := s.H.DB.Where("pool_id = ?", req.PoolId).First(&pool).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find pool with ID %s: %v", req.PoolId, err)
	}

	if pool.UserID != userID {
		return nil, status.Errorf(codes.PermissionDenied, "User is not the owner of the pool with ID %s", req.PoolId)
	}

	opts := galileo.AccountsAndCardsApiPostModifystatusOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(pool.PRN),
		Type_:               optional.NewInt32(16),
	}

	response, _, err := s.GalileoClient.AccountsAndCardsApi.PostModifystatus(ctx, &opts)

	if response.Status != "Success" {
		return nil, status.Errorf(codes.Internal, "Could not create account in Galileo.")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not create pool in Galileo.")
	}

	tx := s.H.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&pool).Update("status", "deleted").Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Could not delete pool: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Transaction commit error: %v", err)
	}

	return &pb.DeletePoolResponse{
		PoolId: req.PoolId,
	}, nil

}

func (s *Server) GetUserPools(ctx context.Context, req *pb.GetUserPoolsRequest) (*pb.GetUserPoolsResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	var pools []models.Pool

	err := s.H.DB.Joins("JOIN participants on participants.pool_id = pools.pool_id").
		Where("participants.user_id = ?", userId).
		Find(&pools).Error

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not find pools for user with ID %s: %v", userId, err)
	}

	var pbPools []*pb.Pool
	for _, pool := range pools {
		pbPools = append(pbPools, utils.ConvertToPBPool(pool))
	}

	return &pb.GetUserPoolsResponse{
		Pools: pbPools,
	}, nil
}

func (s *Server) GetPool(ctx context.Context, req *pb.GetPoolRequest) (*pb.GetPoolResponse, error) {
	var (
		pool         models.Pool
		participants []models.Participant
		transactions []models.Transaction
	)

	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	var participant models.Participant
	if err := s.H.DB.Where("pool_id = ? AND user_id = ?", req.PoolId, userId).First(&participant).Error; err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "User is not a participant in the pool with ID %s", req.PoolId)
	}

	tx := s.H.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Where("pool_id = ?", req.PoolId).First(&pool).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "Could not find pool with ID %s: %v", req.PoolId, err)
	}

	err := tx.Where("pool_id = ?", pool.PoolID).Find(&participants).Error
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Could not find participants for pool with ID %s: %v", req.PoolId, err)
	}

	err = tx.Where("from_account_id = ? OR to_account_id = ?", pool.PRN, pool.PRN).Find(&transactions).Error
	if err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Could not find transactions for pool with ID %s: %v", req.PoolId, err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Transaction commit error: %v", err)
	}

	var pbParticipants []*pb.Participant
	for _, participant := range participants {
		pbParticipants = append(pbParticipants, utils.ConvertToPBParticipant(participant))
	}

	var pbTransactions []*pb.Transaction
	for _, transaction := range transactions {
		pbTransactions = append(pbTransactions, utils.ConvertToPBTransaction(transaction))
	}

	return &pb.GetPoolResponse{
		PoolId:       pool.PoolID.String(),
		Poolname:     pool.PoolName,
		Description:  pool.Description,
		BalanceGoal:  int32(pool.BalanceGoal),
		Participants: pbParticipants,
		Transactions: pbTransactions,
	}, nil
}

func (s *Server) RemovePoolParticipant(ctx context.Context, req *pb.RemovePoolParticipantRequest) (*pb.RemovePoolParticipantResponse, error) {
	userID, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	var pool models.Pool
	if err := s.H.DB.Where("pool_id = ?", req.PoolId).First(&pool).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find pool with ID %s: %v", req.PoolId, err)
	}

	if pool.UserID != userID {
		return nil, status.Errorf(codes.PermissionDenied, "User is not the owner of the pool with ID %s", req.PoolId)
	}

	var participant models.Participant
	if err := s.H.DB.Where("pool_id = ? AND user_id = ?", req.PoolId, req.ParticipantId).First(&participant).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find participant with ID %s in pool with ID %s: %v", req.ParticipantId, req.PoolId, err)
	}

	tx := s.H.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("participant_id = ?", participant.ParticipantID).Delete(&participant).Error; err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Could not remove participant with ID %s from pool with ID %s: %v", req.ParticipantId, req.PoolId, err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Transaction commit error: %v", err)
	}

	return &pb.RemovePoolParticipantResponse{
		PoolId:        req.PoolId,
		ParticipantId: req.ParticipantId,
	}, nil
}

func (s *Server) UpdatePool(ctx context.Context, req *pb.UpdatePoolRequest) (*pb.UpdatePoolResponse, error) {
	userID, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	var pool models.Pool
	if err := s.H.DB.Where("pool_id = ?", req.PoolId).First(&pool).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find pool with ID %s: %v", req.PoolId, err)
	}

	if pool.UserID != userID {
		return nil, status.Errorf(codes.PermissionDenied, "User is not the owner of the pool with ID %s", req.PoolId)
	}

	if req.Description != "" {
		pool.Description = req.Description
	}

	if req.Poolname != "" {
		pool.PoolName = req.Poolname
	}

	if req.BalanceGoal != 0 {
		pool.BalanceGoal = float64(req.BalanceGoal)
	}

	if err := s.H.DB.Save(&pool).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Could not update pool: %v", err)
	}

	return &pb.UpdatePoolResponse{
		Poolname:    pool.PoolName,
		Description: pool.Description,
		BalanceGoal: int32(pool.BalanceGoal),
	}, nil
}

func (s *Server) IsUserOwnerOfPool(ctx context.Context, req *pb.IsUserOwnerOfPoolRequest) (*pb.IsUserOwnerOfPoolResponse, error) {
	var pool models.Pool
	if err := s.H.DB.Where("pool_id = ?", req.PoolId).First(&pool).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find pool with ID %s: %v", req.PoolId, err)
	}

	if pool.UserID != req.UserId {
		return &pb.IsUserOwnerOfPoolResponse{
			IsOwner: false,
		}, nil
	}

	return &pb.IsUserOwnerOfPoolResponse{
		IsOwner: true,
	}, nil
}

func (s *Server) GetPoolBalance(ctx context.Context, req *pb.GetPoolBalanceRequest) (*pb.GetPoolBalanceResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid. - ")
	}

	isMember, err := s.H.IsMemberOfPool(userId, req.PoolId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error checking if user is member of pool. - %v", err)
	}

	if !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "User is not member of pool. - ")
	}

	opts := galileo.AccountsAndCardsApiPostGetbalanceOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(req.PoolId),
	}

	response, _, err := s.GalileoClient.AccountsAndCardsApi.PostGetbalance(ctx, &opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting pool balance. - %v", err)
	}

	return &pb.GetPoolBalanceResponse{
		Balance: response.ResponseData.Balance,
	}, nil
}

func (s *Server) GetPoolOverview(ctx context.Context, req *pb.GetPoolOverviewRequest) (*pb.GetPoolOverviewResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid. - ")
	}

	isMember, err := s.H.IsMemberOfPool(userId, req.PoolId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error checking if user is member of pool. - %v", err)
	}

	if !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "User is not member of pool. - ")
	}

	var pool models.Pool
	if err := s.H.DB.Where("pool_id = ?", req.PoolId).First(&pool).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not find pool with ID %s: %v", req.PoolId, err)
	}

	opts := galileo.TransactionsApiPostGetaccountoverviewOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(pool.PRN),
	}

	response, _, err := s.GalileoClient.TransactionsApi.PostGetaccountoverview(ctx, &opts)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting account overview. - %v", err)
	}

	pbTransactions := utils.ConvertTransactions(response.ResponseData.Transactions)

	return &pb.GetPoolOverviewResponse{
		Balance:          response.ResponseData.Balance,
		Transactions:     pbTransactions,
		TransactionCount: response.ResponseData.TransactionCount,
	}, nil
}
