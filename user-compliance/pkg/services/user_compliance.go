package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	galileo "github.com/Sharefunds/galileo-client"
	"github.com/Sharefunds/user-compliance/pkg/config"
	"github.com/Sharefunds/user-compliance/pkg/db"
	"github.com/Sharefunds/user-compliance/pkg/models"
	"github.com/Sharefunds/user-compliance/pkg/pb"
	"github.com/Sharefunds/user-compliance/pkg/utils"
	"github.com/antihax/optional"
	"github.com/google/uuid"
	plaid "github.com/plaid/plaid-go/v20/plaid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type Server struct {
	H             db.Handler
	GalileoClient *galileo.APIClient
	PlaidClient   *plaid.APIClient
	C             config.Config
	pb.UnimplementedUserComplianceServiceServer
}

// #region User Compliance - Account Operations
func (s *Server) CreateUserRecord(ctx context.Context, req *pb.CreateUserRecordRequest) (*pb.CreateUserRecordResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	newUser := models.User{
		UID:      userId,
		Username: userId,
	}

	if err := s.H.DB.Create(&newUser).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return &pb.CreateUserRecordResponse{
			Message: "Error creating user.",
		}, err
	}

	return &pb.CreateUserRecordResponse{
		Message: "User record created successfully.",
	}, nil

}

func (s *Server) CreateUserAccount(ctx context.Context, req *pb.CreateUserAccountRequest) (*pb.CreateUserAccountResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	log.Println("Starting Transaction for Galileo account creation and user-compliance DB record creation.")
	tx := s.H.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var newUser models.User

	if err := tx.Where("UID = ?", userId).First(&newUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Error fetching user: %v", err)
	}

	// TODO: Fetch user's personal info and address from Plaid API and store in DB and populate in Galileo request

	opts := galileo.AccountsAndCardsApiPostCreateaccountOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		ProdId:              optional.NewInt32(s.C.GalileoProductId),
		TransactionId:       optional.NewString(uuid.New().String()),
	}

	response, _, err := s.GalileoClient.AccountsAndCardsApi.PostCreateaccount(ctx, &opts)

	if response.Status != "Success" {
		log.Printf("Error creating account in Galileo: %v", response.ResponseData)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Error creating account in Galileo: %v", response.ResponseData)
	}

	if err != nil {
		fmt.Println("Error creating account in Galileo.", err)
		tx.Rollback()
		return nil, err
	}

	newUser.PRN = response.ResponseData[0].PmtRefNo
	newUser.Username = req.Username

	//Modify account status to be active

	opts2 := galileo.AccountsAndCardsApiPostModifystatusOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(newUser.PRN),
		Type_:               optional.NewInt32(11),
	}

	response2, _, err := s.GalileoClient.AccountsAndCardsApi.PostModifystatus(ctx, &opts2)

	if response2.Status != "Success" {
		log.Printf("Error modifying account status in Galileo: %v", response2.ResponseData)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Error modifying account status in Galileo: %v", response2.ResponseData)
	}

	if err != nil {
		fmt.Println("Error modifying account status in Galileo.", err)
		tx.Rollback()
		return nil, err
	}

	if err := tx.Save(&newUser).Error; err != nil {
		log.Printf("Error saving user: %v", err)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Error saving user: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	log.Println("Transaction completed successfully.")

	fmt.Println(response.ResponseData)

	return &pb.CreateUserAccountResponse{
		Message: "User account created successfully.",
		GalileoAccountDetails: &pb.GalileoAccountDetails{
			CardId:           response.ResponseData[0].CardId,
			CardNumber:       response.ResponseData[0].CardNumber,
			CardSecurityCode: response.ResponseData[0].CardSecurityCode,
		},
		GalileoAccountStatus: "Active",
	}, nil
}

func (s *Server) GetUserAccount(ctx context.Context, req *pb.GetUserAccountRequest) (*pb.GetUserAccountResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	fmt.Println(userId)

	var user models.User
	result := s.H.DB.Where("UID = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		return nil, status.Errorf(codes.Internal, "Error fetching user: %v", result.Error)
	}

	return &pb.GetUserAccountResponse{
		Username: user.Username,
		PersonalInfo: &pb.PersonalInfo{
			FirstName:   user.PersonalInfo.FirstName,
			LastName:    user.PersonalInfo.LastName,
			DateOfBirth: timestamppb.New(user.PersonalInfo.DateOfBirth),
			Email:       user.PersonalInfo.Email,
			PhoneNumber: user.PersonalInfo.PhoneNumber,
		},
		Address: &pb.Address{
			Address1:    user.Address.Address1,
			Address2:    user.Address.Address2,
			City:        user.Address.City,
			State:       user.Address.State,
			PostalCode:  user.Address.PostalCode,
			CountryCode: user.Address.CountryCode,
		},
	}, nil
}

func (s *Server) UpdateUserAccount(ctx context.Context, req *pb.UpdateUserAccountRequest) (*pb.UpdateUserAccountResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	tx := s.H.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	if err := tx.Where("UID = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Error fetching user: %v", err)
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.PersonalInfo != nil {
		if req.PersonalInfo.FirstName != "" {
			user.PersonalInfo.FirstName = req.PersonalInfo.FirstName
		}
		if req.PersonalInfo.LastName != "" {
			user.PersonalInfo.LastName = req.PersonalInfo.LastName
		}
		if req.PersonalInfo.Email != "" {
			user.PersonalInfo.Email = req.PersonalInfo.Email
		}
		if req.PersonalInfo.PhoneNumber != "" {
			user.PersonalInfo.PhoneNumber = req.PersonalInfo.PhoneNumber
		}
		if req.PersonalInfo.DateOfBirth != nil {
			user.PersonalInfo.DateOfBirth = req.PersonalInfo.DateOfBirth.AsTime()
		}
	}

	if req.Address != nil {
		if req.Address.Address1 != "" {
			user.Address.Address1 = req.Address.Address1
		}
		if req.Address.Address2 != "" {
			user.Address.Address2 = req.Address.Address2
		}
		if req.Address.City != "" {
			user.Address.City = req.Address.City
		}
		if req.Address.State != "" {
			user.Address.State = req.Address.State
		}
		if req.Address.PostalCode != "" {
			user.Address.PostalCode = req.Address.PostalCode
		}
		if req.Address.CountryCode != "" {
			user.Address.CountryCode = req.Address.CountryCode
		}
	}

	dobTime := user.PersonalInfo.DateOfBirth.Format("2006-01-02")

	opts := galileo.AccountsAndCardsApiPostUpdateaccountOpts{
		ResponseContentType: optional.NewString("json"),
		ApiLogin:            optional.NewString(s.C.GalileoLogin),
		ApiTransKey:         optional.NewString(s.C.GalileoTranskey),
		ProviderId:          optional.NewInt32(s.C.GalileoProviderId),
		TransactionId:       optional.NewString(uuid.New().String()),
		AccountNo:           optional.NewString(user.PRN),
		FirstName:           optional.NewString(user.PersonalInfo.FirstName),
		LastName:            optional.NewString(user.PersonalInfo.LastName),
		DateOfBirth:         optional.NewString(dobTime),
		Address1:            optional.NewString(user.Address.Address1),
		Address2:            optional.NewString(user.Address.Address2),
		City:                optional.NewString(user.Address.City),
		State:               optional.NewString(user.Address.State),
		PostalCode:          optional.NewString(user.Address.PostalCode),
		CountryCode:         optional.NewString("840"),
		Email:               optional.NewString(user.PersonalInfo.Email),
		PrimaryPhone:        optional.NewString(user.PersonalInfo.PhoneNumber),
	}

	response, _, err := s.GalileoClient.AccountsAndCardsApi.PostUpdateaccount(ctx, &opts)

	if response.Status != "Success" {
		log.Printf("Error updating account in Galileo: %v", response.ResponseData)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Error updating account in Galileo: %v", response.ResponseData)
	}

	if err != nil {
		log.Printf("Error updating account in Galileo: %v", err)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Error updating account in Galileo: %v", err)
	}

	if err := tx.Save(&user).Error; err != nil {
		log.Printf("Error updating user: %v", err)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "Error updating user: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Error committing transaction: %v", err)
	}

	return &pb.UpdateUserAccountResponse{
		Message: "User profile updated successfully.",
		GalileoAccountInfo: &pb.GalileoAccountInfo{
			PersonalInfo: &pb.PersonalInfo{
				FirstName: response.ResponseData.CustomerProfile.FirstName,
				LastName:  response.ResponseData.CustomerProfile.LastName,
				Email:     response.ResponseData.CustomerProfile.Email,
			},
			Address: &pb.Address{
				Address1:    response.ResponseData.CustomerProfile.Address1,
				Address2:    response.ResponseData.CustomerProfile.Address2,
				City:        response.ResponseData.CustomerProfile.City,
				State:       response.ResponseData.CustomerProfile.State,
				PostalCode:  response.ResponseData.CustomerProfile.PostalCode,
				CountryCode: response.ResponseData.CustomerProfile.CountryCode,
			},
		},
	}, nil
}

func (s *Server) GetUserPRN(ctx context.Context, req *pb.GetUserPRNRequest) (*pb.GetUserPRNResponse, error) {
	userId := req.GetUserId()

	var user models.User
	result := s.H.DB.Where("UID = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		return nil, status.Errorf(codes.Internal, "Error fetching user: %v", result.Error)
	}

	return &pb.GetUserPRNResponse{
		Prn: user.PRN,
	}, nil
}

//#endregion

//#region Plaid API

func (s *Server) CreatePlaidLinkToken(ctx context.Context, req *pb.CreatePlaidLinkTokenRequest) (*pb.CreatePlaidLinkTokenResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: userId,
	}

	request := plaid.NewLinkTokenCreateRequest(
		"PoolParty",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		user,
	)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH})
	request.SetRequiredIfSupportedProducts([]plaid.Products{plaid.PRODUCTS_IDENTITY})

	linkTokenCreateResp, _, err := s.PlaidClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		panic(err)
	}
	linkToken := linkTokenCreateResp.GetLinkToken()

	return &pb.CreatePlaidLinkTokenResponse{
		LinkToken: linkToken,
	}, nil
}

func (s *Server) CreatePlaidAccessToken(ctx context.Context, req *pb.CreatePlaidAccessTokenRequest) (*pb.CreatePlaidAccessTokenResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	exchangePublicTokenReq := plaid.NewItemPublicTokenExchangeRequest(req.PublicToken)
	exchangePublicTokenResp, _, err := s.PlaidClient.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*exchangePublicTokenReq,
	).Execute()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error exchanging public token: %v", err)
	}
	accessToken := exchangePublicTokenResp.GetAccessToken()

	if err := s.storeAccessTokenData(userId, accessToken, req.Metadata); err != nil {
		return nil, status.Errorf(codes.Internal, "Error storing access token data: %v", err)
	}

	processorToken, err := s.createProcessorToken(ctx, accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating processor token: %v", err)
	}

	fmt.Println(processorToken)

	return &pb.CreatePlaidAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *Server) createProcessorToken(ctx context.Context, accessToken string) (string, error) {
	accountsResp, _, err := s.PlaidClient.PlaidApi.AccountsGet(ctx).AccountsGetRequest(plaid.AccountsGetRequest{
		AccessToken: accessToken,
	}).Execute()

	if err != nil {
		return "", err
	}

	accountID := accountsResp.GetAccounts()[0].GetAccountId()

	request := plaid.NewProcessorTokenCreateRequest(accessToken, accountID, "galileo")

	stripeTokenResp, _, err := s.PlaidClient.PlaidApi.ProcessorTokenCreate(ctx).ProcessorTokenCreateRequest(
		*request,
	).Execute()

	if err != nil {
		return "", err
	}

	return stripeTokenResp.GetProcessorToken(), nil
}

func (s *Server) storeAccessTokenData(userID string, accessToken string, metadata *pb.PlaidMetadata) error {

	var bankAccount models.BankAccount
	result := s.H.DB.Where(models.BankAccount{UserID: userID, InstitutionName: metadata.Institution.Name}).FirstOrCreate(&bankAccount, models.BankAccount{
		UserID:          userID,
		InstitutionName: metadata.Institution.GetName(),
		AccountName:     metadata.GetAccounts()[0].Name,
		AccountType:     metadata.GetAccounts()[0].Type,
		Mask:            metadata.GetAccounts()[0].Mask,
	})

	if result.Error != nil {
		return result.Error
	}

	plaidItem := models.PlaidItem{
		UserID:        userID,
		BankAccountID: bankAccount.ID,
		AccessToken:   accessToken,
		InstitutionID: metadata.Institution.GetId(),
	}

	result = s.H.DB.Create(&plaidItem)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Server) fetchAccessToken(ctx context.Context, userID string, institutionID string) (string, error) {
	var plaidItem models.PlaidItem
	result := s.H.DB.Where("user_id = ? AND institution_id = ?", userID, institutionID).First(&plaidItem)

	if result.Error != nil {
		return "", result.Error
	}

	return plaidItem.AccessToken, nil
}

func (s *Server) fetchAccessTokenData(ctx context.Context, userID string) ([]*models.PlaidItem, error) {
	var plaidItems []*models.PlaidItem
	result := s.H.DB.Where("user_id = ?", userID).Find(&plaidItems)

	if result.Error != nil {
		return nil, result.Error
	}

	return plaidItems, nil
}

func (s *Server) CreatePlaidIDVToken(ctx context.Context, req *pb.CreatePlaidIDVTokenRequest) (*pb.CreatePlaidIDVTokenResponse, error) {
	userId, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: userId,
	}

	identityVerification := plaid.LinkTokenCreateRequestIdentityVerification{
		TemplateId: "idvtmp_1vPccox1Ln2tgG",
	}

	request := plaid.NewLinkTokenCreateRequest(
		"PoolParty",
		"en",
		[]plaid.CountryCode{plaid.COUNTRYCODE_US},
		user,
	)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_IDENTITY_VERIFICATION})
	request.SetIdentityVerification(identityVerification)

	linkTokenCreateResp, _, err := s.PlaidClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		panic(err)
	}

	linkToken := linkTokenCreateResp.GetLinkToken()

	return &pb.CreatePlaidIDVTokenResponse{
		IdvToken: linkToken,
	}, nil

}

func (s *Server) SetIdentityVerified(ctx context.Context, req *pb.SetIdentityVerifiedRequest) (*pb.SetIdentityVerifiedResponse, error) {
	userID, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	var user models.User
	result := s.H.DB.Where("UID = ?", userID).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		return nil, status.Errorf(codes.Internal, "Error fetching user: %v", result.Error)
	}

	// TODO: Retrieve identity status from Plaid API

	user.IsVerified = true
	user.IdentityToken = req.LinkSessionId

	result = s.H.DB.Save(&user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "Error saving user: %v", result.Error)
	}

	return &pb.SetIdentityVerifiedResponse{
		IsVerified: true,
	}, nil
}

func (s *Server) GetAccountStatus(ctx context.Context, req *pb.GetAccountStatusRequest) (*pb.GetAccountStatusResponse, error) {
	userID, ok := ctx.Value(utils.UserIdKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not found in context or invalid.")
	}

	var user models.User
	result := s.H.DB.Where("UID = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "User not found.")
		}
		return nil, status.Errorf(codes.Internal, "Error fetching user: %v", result.Error)
	}

	return &pb.GetAccountStatusResponse{
		GalileoAccountCreated: user.PRN != "",
		IsVerified:            user.IsVerified,
	}, nil

}

//#endregion
