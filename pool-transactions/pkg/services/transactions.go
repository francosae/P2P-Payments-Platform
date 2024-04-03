package services

import (
	"context"

	galileo "github.com/Sharefunds/galileo-client"
	"github.com/Sharefunds/pool-transactions/pkg/utils"
	"github.com/antihax/optional"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sharefunds/pool-transactions/pkg/pb"
)

func (s *Server) GetUserBalance(ctx context.Context, req *pb.GetUserBalanceRequest) (*pb.GetUserBalanceResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not provided")
	}

	userPrn, err := s.UserComplianceClient.GetUserPRN(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get PRN from User Compliance Service.")
	}

	opts := galileo.AccountsAndCardsApiPostGetbalanceOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(userPrn.Prn),
	}

	response, _, err := s.GalileoClient.AccountsAndCardsApi.PostGetbalance(ctx, &opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting pool balance. - %v", err)
	}

	return &pb.GetUserBalanceResponse{
		Balance: response.ResponseData.Balance,
	}, nil

}

func (s *Server) SendPaymentToUser(ctx context.Context, req *pb.SendPaymentToUserRequest) (*pb.SendPaymentToUserResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not provided")
	}

	senderUserPrn, err := s.UserComplianceClient.GetUserPRN(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get PRN from User Compliance Service.")
	}

	receiverUserPrn, err := s.UserComplianceClient.GetUserPRN(ctx, req.ReceiverUserId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get PRN from User Compliance Service.")
	}

	opts := galileo.TransactionsApiPostCreateaccounttransferOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(senderUserPrn.Prn),
		Amount:              optional.NewFloat32(req.Amount),
		TransferToAccountNo: optional.NewString(receiverUserPrn.Prn),
	}

	response, _, err := s.GalileoClient.TransactionsApi.PostCreateaccounttransfer(ctx, &opts)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error sending payment. - %v", err)
	}

	return &pb.SendPaymentToUserResponse{
		OldBalance: response.ResponseData.OldBalance,
		NewBalance: response.ResponseData.NewBalance,
	}, nil

}

func (s *Server) SendPaymentToPool(ctx context.Context, req *pb.SendPaymentToPoolRequest) (*pb.SendPaymentToPoolResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not provided")
	}

	userPrn, err := s.UserComplianceClient.GetUserPRN(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get PRN from User Compliance Service.")
	}

	opts := galileo.TransactionsApiPostCreateaccounttransferOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(userPrn.Prn),
		Amount:              optional.NewFloat32(req.Amount),
		TransferToAccountNo: optional.NewString(req.PoolId),
	}

	response, _, err := s.GalileoClient.TransactionsApi.PostCreateaccounttransfer(ctx, &opts)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error sending payment. - %v", err)
	}

	return &pb.SendPaymentToPoolResponse{
		OldBalance: response.ResponseData.OldBalance,
		NewBalance: response.ResponseData.NewBalance,
	}, nil

}

func (s *Server) GetUserAccountOverview(ctx context.Context, req *pb.GetUserAccountOverviewRequest) (*pb.GetUserAccountOverviewResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not provided")
	}

	userPrn, err := s.UserComplianceClient.GetUserPRN(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get PRN from User Compliance Service.")
	}

	opts := galileo.TransactionsApiPostGetaccountoverviewOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(userPrn.Prn),
	}

	response, _, err := s.GalileoClient.TransactionsApi.PostGetaccountoverview(ctx, &opts)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting account overview. - %v", err)
	}

	pbTransactions := utils.ConvertTransactions(response.ResponseData.Transactions)

	return &pb.GetUserAccountOverviewResponse{
		Balance:          response.ResponseData.Balance,
		Transactions:     pbTransactions,
		TransactionCount: response.ResponseData.TransactionCount,
	}, nil
}
