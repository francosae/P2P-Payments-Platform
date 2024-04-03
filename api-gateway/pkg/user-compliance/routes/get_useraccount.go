package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
)

type GetUserProfileRequestBody struct {
}

// @Summary Retrieve a user's account info stored on the user-compliance service DB.
// @Description Retrieves a user's account info stored on the user-compliance service DB. Does not make request to Galileo for user account info.
// @Tags User & Compliance Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UID path string true "User ID"
// @Success 200 {object} pb.GetUserAccountResponse
// @ID get_user_account
// @Router /user-compliance/users/get-account/{UID} [get]
func GetUserAccount(ctx *gin.Context, c pb.UserComplianceServiceClient) {

	grpcCtx := ctx.Request.Context()

	res, err := c.GetUserAccount(grpcCtx, &pb.GetUserAccountRequest{})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
