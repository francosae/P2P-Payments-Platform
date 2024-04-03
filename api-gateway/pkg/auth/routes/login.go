package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Sharefunds/api-gateway/pkg/auth/pb"
)

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Login a user
// @Description Login a user, validates client token with firebase and returns the client token
// @Tags Authentication Service
// @Accept json
// @Produce json
// @Param user body LoginRequestBody true "User to login"
// @Success 200 {object} pb.LoginResponse
// @Router /auth/login [post]
func Login(ctx *gin.Context, c pb.AuthServiceClient) {
	b := LoginRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.Login(context.Background(), &pb.LoginRequest{
		Email:    b.Email,
		Password: b.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(int(res.Status), &res)
}
