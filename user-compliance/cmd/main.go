package main

import (
	"log"
	"net"

	galileo "github.com/Sharefunds/galileo-client"
	"github.com/Sharefunds/user-compliance/pkg/config"
	"github.com/Sharefunds/user-compliance/pkg/db"
	"github.com/Sharefunds/user-compliance/pkg/pb"
	"github.com/Sharefunds/user-compliance/pkg/services"
	"github.com/Sharefunds/user-compliance/pkg/utils"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	plaid "github.com/plaid/plaid-go/v20/plaid"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	galileoConfig := galileo.NewConfiguration()
	galileoConfig.BasePath = c.GalileoUrl
	galileoClient := galileo.NewAPIClient(galileoConfig)

	plaidConfig := plaid.NewConfiguration()
	plaidConfig.AddDefaultHeader("PLAID-CLIENT-ID", c.PlaidClientId)
	plaidConfig.AddDefaultHeader("PLAID-SECRET", c.PlaidSandboxId)
	plaidConfig.UseEnvironment(plaid.Sandbox)
	plaidClient := plaid.NewAPIClient(plaidConfig)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	log.Println("Sharefunds User & Compliance Service listening on port", c.Port)
	s := services.Server{H: h, GalileoClient: galileoClient, PlaidClient: plaidClient, C: c}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			utils.UserIdInterceptor,
			utils.LoggingInterceptor(),
		)),
	)
	pb.RegisterUserComplianceServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
