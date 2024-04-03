package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

type RemovePoolParticipantRequestBody struct {
	UserId string `json:"UserId"`
}

// @Summary Remove a participant from a pool.
// @Description Removes a participant from a pool, using the poolId and participantId, user must be the owner of the pool.
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Pool body RemovePoolParticipantRequestBody true "User ID to remove from pool"
// @Param PoolID path string true "Pool ID"
// @Success 200 {object} pb.RemovePoolParticipantResponse
// @ID remove_pool_participant
// @Router /pool-transactions/pools/{PoolID}/remove-participant [put]
func RemovePoolParticipant(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	body := RemovePoolParticipantRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	poolID := ctx.Param("PoolID")

	grpcCtx := ctx.Request.Context()

	res, err := c.RemovePoolParticipant(grpcCtx, &pb.RemovePoolParticipantRequest{
		PoolId:        poolID,
		ParticipantId: body.UserId,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
