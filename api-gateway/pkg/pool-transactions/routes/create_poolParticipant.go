package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Creates ( accepts ) a pool participant , from an invitation.
// @Description Accepts pool invitation, using the poolId.
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PoolID path string true "Pool ID"
// @Success 200 {object} pb.CreatePoolParticipantResponse
// @ID create_pool_participant
// @Router /pool-transactions/pools/{PoolID}/accept-invite [post]
func CreatePoolParticipant(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	poolID := ctx.Param("PoolID")

	grpcCtx := ctx.Request.Context()

	res, err := c.CreatePoolParticipant(grpcCtx, &pb.CreatePoolParticipantRequest{
		PoolId: poolID,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
