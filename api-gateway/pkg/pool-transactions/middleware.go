package pooltransactions

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type PoolTransactionsMiddlewareConfig struct {
	svc *ServiceClient
}

func InitPoolTransactionsMiddleware(svc *ServiceClient) PoolTransactionsMiddlewareConfig {
	return PoolTransactionsMiddlewareConfig{svc}
}

func (c *PoolTransactionsMiddlewareConfig) PoolOwnerRequired(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - User ID not found"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - User ID is not a string"})
		return
	}

	poolID := ctx.Param("PoolID")
	md := metadata.Pairs("user_id", userIDStr)
	grpcCtx := metadata.NewOutgoingContext(ctx.Request.Context(), md)

	res, err := c.svc.Client.IsUserOwnerOfPool(grpcCtx, &pb.IsUserOwnerOfPoolRequest{
		PoolId: poolID,
		UserId: userID.(string),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !res.IsOwner {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden - User is not the owner of the pool"})
		return
	}

	ctx.Next()
}
