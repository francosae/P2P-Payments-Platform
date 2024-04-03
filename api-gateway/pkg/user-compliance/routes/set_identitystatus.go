package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
)

type SetIdentityStatusRequestBody struct {
	LinkSessionId string `json:"link_session_id"`
}

// @Summary Sets a user's identity status, which is stored on the user-compliance service DB.
// @Description Sets a boolean value indicating whether the user has completed the identity verification flow.
// @Tags User & Compliance Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UID path string true "User ID"
// @Param LinkSessionId body SetIdentityStatusRequestBody true "Session ID of the link token"
// @Success 200 {object} pb.SetIdentityVerifiedResponse
// @ID set_identity_status
// @Router /user-compliance/users/set-verification-status/{UID} [post]
func SetIdentityStatus(ctx *gin.Context, c pb.UserComplianceServiceClient) {
	body := SetIdentityStatusRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.SetIdentityVerified(grpcCtx, &pb.SetIdentityVerifiedRequest{
		LinkSessionId: body.LinkSessionId,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
