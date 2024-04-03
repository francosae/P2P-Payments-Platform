package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Gets a user's account status, which is stored on the user-compliance service DB.
// @Description Gets a boolean value indicating whether the user has completed the identity verification flow, and galileo account creation.
// @Tags User & Compliance Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UID path string true "User ID"
// @Success 200 {object} pb.GetAccountStatusResponse
// @ID get_account_status
// @Router /user-compliance/users/get-account-status/{UID} [get]
func GetAccountStatus(ctx *gin.Context, c pb.UserComplianceServiceClient) {

	grpcCtx := ctx.Request.Context()

	res, err := c.GetAccountStatus(grpcCtx, &pb.GetAccountStatusRequest{})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
