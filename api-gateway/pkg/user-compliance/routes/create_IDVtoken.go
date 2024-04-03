package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Creates an Identity Verification Link Token from the Plaid API.
// @Description Creates a link token for the Identity Verification flow, which allows the user to verify their identity.
// @Tags User & Compliance Service - Plaid
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UID path string true "User ID"
// @Success 200 {object} pb.CreatePlaidIDVTokenResponse
// @ID create_IDV_token
// @Router /user-compliance/plaid/create-idvtoken/{UID} [post]
func CreateIDVToken(ctx *gin.Context, c pb.UserComplianceServiceClient) {

	grpcCtx := ctx.Request.Context()

	res, err := c.CreatePlaidIDVToken(grpcCtx, &pb.CreatePlaidIDVTokenRequest{})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
