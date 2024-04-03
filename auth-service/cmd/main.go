package main

import (
	"log"
	"net"

	"github.com/Sharefunds/auth-service/pkg/config"
	"github.com/Sharefunds/auth-service/pkg/db"
	"github.com/Sharefunds/auth-service/pkg/firebase"
	"github.com/Sharefunds/auth-service/pkg/galileo"
	"github.com/Sharefunds/auth-service/pkg/pb"
	"github.com/Sharefunds/auth-service/pkg/services"

	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	galileoClient := galileo.InitGalileoClient(
		c.GalileoUrl,
		c.GalileoLogin,
		c.GalileoTranskey,
		c.GalileoProviderId,
		c.GalileoProductId,
	)

	firebaseClient, err := firebase.InitializeAppWithServiceAccount()

	if err != nil {
		log.Fatalln("Failed to initialize Firebase", err)
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	log.Println("Sharefunds Authentication Service listening on port", c.Port)
	s := services.Server{H: h, GalileoClient: galileoClient, FirebaseClient: firebaseClient}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
