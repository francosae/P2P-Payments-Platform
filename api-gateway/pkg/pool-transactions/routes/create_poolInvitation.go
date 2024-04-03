package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

type CreatePoolInvitationRequestBody struct {
	UserId string `json:"UserId"`
	// Role   string `json:"Role"`
}

// @Summary Create Pool Invitation, creates a pool invitation for a user.
// @Description Creates a pool invitation for a user.
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Pool body CreatePoolInvitationRequestBody true "User ID to invite"
// @Param PoolID path string true "Pool ID"
// @Success 200 {object} pb.CreatePoolInvitationResponse
// @ID create_pool_invitation
// @Router /pool-transactions/pools/{PoolID}/create-invite [post]
func CreatePoolInvitation(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	poolID := ctx.Param("PoolID")

	body := CreatePoolInvitationRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.CreatePoolInvitation(grpcCtx, &pb.CreatePoolInvitationRequest{
		PoolId: poolID,
		UserId: body.UserId,
		Role:   "member",
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
