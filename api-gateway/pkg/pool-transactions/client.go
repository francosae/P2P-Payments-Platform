package pooltransactions

import (
	"fmt"

	"github.com/Sharefunds/api-gateway/pkg/config"
	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceClient struct {
	Client pb.PoolTransactionsServiceClient
}

func InitServiceClient(c *config.Config) pb.PoolTransactionsServiceClient {
	cc, err := grpc.Dial(c.PoolTransactionSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		fmt.Println("Could not connect to:", err)
	}

	return pb.NewPoolTransactionsServiceClient(cc)
}
