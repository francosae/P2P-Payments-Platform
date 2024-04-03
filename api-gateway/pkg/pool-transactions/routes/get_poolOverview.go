package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Get Pool Overview.
// @Description Fetches the overview of a pool, including pool's balance, transactions, and transaction count.
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PoolID path string true "Pool ID"
// @Success 200 {object} pb.GetPoolOverviewResponse
// @ID get_pool_overview
// @Router /pool-transactions/pools/overview/{PoolID} [get]
func GetPoolOverview(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	grpcCtx := ctx.Request.Context()

	res, err := c.GetPoolOverview(grpcCtx, &pb.GetPoolOverviewRequest{
		PoolId: ctx.Param("PoolID"),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
