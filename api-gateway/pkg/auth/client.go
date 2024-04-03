package auth

import (
	"fmt"

	"github.com/Sharefunds/api-gateway/pkg/auth/pb"
	"github.com/Sharefunds/api-gateway/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client pb.AuthServiceClient
}

func InitServiceClient(c *config.Config) pb.AuthServiceClient {
	// insecure used for dev . For production will need TLS/SSL
	cc, err := grpc.Dial(c.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect to auth service:", err)
	}

	return pb.NewAuthServiceClient(cc)

}
