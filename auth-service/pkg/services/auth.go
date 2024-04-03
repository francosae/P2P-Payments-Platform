package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/Sharefunds/auth-service/pkg/db"
	"github.com/Sharefunds/auth-service/pkg/firebase"
	"github.com/Sharefunds/auth-service/pkg/galileo"

	"github.com/Sharefunds/auth-service/pkg/galileo/accounts-and-cards/enrollment"
	"github.com/Sharefunds/auth-service/pkg/models"
	"github.com/Sharefunds/auth-service/pkg/pb"
)

type Server struct {
	H              db.Handler
	GalileoClient  *galileo.GalileoClient
	FirebaseClient *firebase.FirebaseClient
	pb.UnimplementedAuthServiceServer
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return &pb.RegisterResponse{Status: http.StatusConflict, Error: "Email already exists"}, nil
	}

	user.Email = req.Email

	client, err := s.FirebaseClient.App.Auth(ctx)
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failure getting Firebase Auth client",
		}, err
	}

	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(req.Password)

	u, err := client.CreateUser(ctx, params)
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failure creating Firebase user",
		}, err
	}

	user.UID = u.UID
	user.Email = u.Email

	if result := s.H.DB.Create(&user); result.Error != nil {
		return &pb.RegisterResponse{Status: http.StatusInternalServerError, Error: "Failed to create user"}, result.Error
	}

	//Galileo Account creation : to be moved to a separate service (user & compliance service)
	accountResponse, error := enrollment.CreateAccount(s.GalileoClient)

	if error != nil {
		return &pb.RegisterResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to Galileo account",
		}, error
	}
	fmt.Println("accountResponse", accountResponse.ResponseData[0])

	createdAccount, err := json.MarshalIndent(accountResponse, "", "    ")
	if err != nil {
		log.Fatalf("Failed to generate pretty JSON: %v", err)
	}
	fmt.Printf("%s\n", createdAccount)
	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	client, err := s.FirebaseClient.App.Auth(ctx)

	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failure getting Firebase Auth client",
		}, err
	}

	decodedToken, err := client.VerifyIDToken(ctx, req.Firebaseidtoken)

	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failure verifying Firebase ID token",
		}, err
	}

	if result := s.H.DB.Where(&models.User{UID: decodedToken.UID}).First(&user); result.Error != nil {
		return &pb.LoginResponse{Status: http.StatusNotFound, Error: "User not found"}, nil
	}

	return &pb.LoginResponse{
		Status: http.StatusOK,
		Token:  req.Firebaseidtoken,
	}, nil

}

func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// var user models.User

	client, err := s.FirebaseClient.App.Auth(ctx)

	if err != nil {
		return &pb.ValidateTokenResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failure getting Firebase Auth client",
		}, err
	}

	decodedToken, err := client.VerifyIDToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failure verifying Firebase ID token",
		}, err
	}

	// TODO: This is temporary. We should have this check added once we incorporate a consistent user creation flow.

	// user.UID = decodedToken.UID

	// if result := s.H.DB.Where(&models.User{UID: decodedToken.UID}).First(&user); result.Error != nil {
	// 	return &pb.ValidateTokenResponse{Status: http.StatusNotFound, Error: "User not found"}, nil
	// }

	return &pb.ValidateTokenResponse{
		Status: http.StatusOK,
		UserId: decodedToken.UID,
	}, nil
}
