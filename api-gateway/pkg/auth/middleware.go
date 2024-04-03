package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/Sharefunds/api-gateway/pkg/auth/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type AuthMiddlewareConfig struct {
	svc *ServiceClient
}

func InitAuthMiddleware(svc *ServiceClient) AuthMiddlewareConfig {
	return AuthMiddlewareConfig{svc}
}

func (c *AuthMiddlewareConfig) AuthRequired(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("Authorization")

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	res, err := c.svc.Client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: token[1],
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("user_id", res.UserId)

	md := metadata.Pairs("user_id", res.UserId)

	newCtx := metadata.NewOutgoingContext(ctx.Request.Context(), md)

	ctx.Request = ctx.Request.WithContext(newCtx)

	ctx.Next()
}

func (c *AuthMiddlewareConfig) UserIdRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, exists := ctx.Get("user_id")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "User ID not found or not valid."})
			return
		}
		ctx.Next()
	}
}
