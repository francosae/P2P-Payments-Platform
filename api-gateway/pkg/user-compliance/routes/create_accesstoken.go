package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
)

type CreateAccessTokenRequestBody struct {
	PublicToken string           `json:"publicToken"`
	Metadata    pb.PlaidMetadata `json:"metadata"`
}

// @Summary Creates an access token from the Plaid API.
// @Description Exchanges the public token for an access token.
// @Tags User & Compliance Service - Plaid
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UID path string true "User ID"
// @Param PublicToken body CreateAccessTokenRequestBody true "Public Token to exchange for access token"
// @Success 200 {object} pb.CreatePlaidAccessTokenResponse
// @ID create_access_token
// @Router /user-compliance/plaid/create-accesstoken/{UID} [post]
func CreateAccessToken(ctx *gin.Context, c pb.UserComplianceServiceClient) {
	body := CreateAccessTokenRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.CreatePlaidAccessToken(grpcCtx, &pb.CreatePlaidAccessTokenRequest{
		PublicToken: body.PublicToken,
		Metadata:    &body.Metadata,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
