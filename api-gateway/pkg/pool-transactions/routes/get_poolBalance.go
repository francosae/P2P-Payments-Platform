package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Get pool balance.
// @Description Fetches the balance of a pool, using the poolId.
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PoolID path string true "Pool ID"
// @Success 200 {object} pb.GetPoolBalanceResponse
// @ID get_pool_balance
// @Router /pool-transactions/pools/balance/{PoolID} [get]
func GetPoolBalance(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	grpcCtx := ctx.Request.Context()

	res, err := c.GetPoolBalance(grpcCtx, &pb.GetPoolBalanceRequest{
		PoolId: ctx.Param("PoolID"),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
