package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Sharefunds/api-gateway/pkg/auth/pb"
)

type RegisterRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Register a new user
// @Description Register a new user, creates an account in Firebase and a user in the database
// @Tags Authentication Service
// @Accept json
// @Produce json
// @Param user body RegisterRequestBody true "User to register"
// @Success 200 {object} pb.RegisterResponse
// @Router /auth/register [post]
func Register(ctx *gin.Context, c pb.AuthServiceClient) {
	body := RegisterRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.Register(context.Background(), &pb.RegisterRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(int(res.Status), &res)
}
