package client

import (
	"context"
	"fmt"

	"github.com/Sharefunds/pool-transactions/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type UserComplianceServiceClient struct {
	Client pb.UserComplianceServiceClient
}

func InitUserComplianceServiceClient(url string) UserComplianceServiceClient {
	cc, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect to:", err)
	}

	c := UserComplianceServiceClient{
		Client: pb.NewUserComplianceServiceClient(cc),
	}

	return c
}

func (c *UserComplianceServiceClient) GetUserPRN(ctx context.Context, userId string) (*pb.GetUserPRNResponse, error) {

	req := &pb.GetUserPRNRequest{
		UserId: userId,
	}

	md := metadata.Pairs("user_id", userId)
	ctx = metadata.NewOutgoingContext(ctx, md)

	return c.Client.GetUserPRN(ctx, req)
}
