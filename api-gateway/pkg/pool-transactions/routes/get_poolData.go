package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Fetch pool data, including participants and transactions.
// @Description Fetches pool data, using the poolId, user must be a participant in the pool.
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PoolID path string true "Pool ID"
// @Success 200 {object} pb.GetPoolResponse
// @ID get_pool_data
// @Router /pool-transactions/pools/{PoolID} [get]
func GetPoolData(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	poolID := ctx.Param("PoolID")

	grpcCtx := ctx.Request.Context()

	res, err := c.GetPool(grpcCtx, &pb.GetPoolRequest{
		PoolId: poolID,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
