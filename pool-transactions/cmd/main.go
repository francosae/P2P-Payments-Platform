package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Sharefunds/pool-transactions/pkg/client"
	"github.com/Sharefunds/pool-transactions/pkg/config"
	"github.com/Sharefunds/pool-transactions/pkg/utils"

	galileo "github.com/Sharefunds/galileo-client"
	"github.com/Sharefunds/pool-transactions/pkg/db"
	"github.com/Sharefunds/pool-transactions/pkg/pb"
	"github.com/Sharefunds/pool-transactions/pkg/services"
	"google.golang.org/grpc"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Error loading config", err)
	}

	dbHandler := db.Init(config.DBUrl)

	usercomplianceClient := client.InitUserComplianceServiceClient(config.UserComplianceSvcUrl)

	galileoConfig := galileo.NewConfiguration()
	galileoConfig.BasePath = config.GalileoUrl
	galileoClient := galileo.NewAPIClient(galileoConfig)

	server := services.Server{
		H:                    dbHandler,
		GalileoClient:        galileoClient,
		UserComplianceClient: usercomplianceClient,
		C:                    config,
	}

	listener, err := net.Listen("tcp", config.Port)
	if err != nil {
		log.Fatalln("Error creating listener", err)
	}

	fmt.Println("Sharefunds Pool & Transactions Service listening on port", config.Port)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(utils.UserIdInterceptor))

	pb.RegisterPoolTransactionsServiceServer(grpcServer, &server)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln("Error creating grpc server", err)
	}
}
