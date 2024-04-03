package usercompliance

import (
	"fmt"

	"github.com/Sharefunds/api-gateway/pkg/config"
	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client pb.UserComplianceServiceClient
}

func InitServiceClient(c *config.Config) pb.UserComplianceServiceClient {
	cc, err := grpc.Dial(c.UserComplianceSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect to:", err)
	}

	return pb.NewUserComplianceServiceClient(cc)
}
