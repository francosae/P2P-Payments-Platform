package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Creates a link token from the Plaid API.
// @Description Creates a link token from the Plaid API.
// @Tags User & Compliance Service - Plaid
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UID path string true "User ID"
// @Success 200 {object} pb.CreatePlaidLinkTokenResponse
// @ID create_link_token
// @Router /user-compliance/plaid/create-linktoken/{UID} [post]
func CreateLinkToken(ctx *gin.Context, c pb.UserComplianceServiceClient) {

	grpcCtx := ctx.Request.Context()

	res, err := c.CreatePlaidLinkToken(grpcCtx, &pb.CreatePlaidLinkTokenRequest{})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
