package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Delete Pool, deletes account in galileo.
// @Description Sets Pool status to Inactive, Deactivates secondary account in galileo.
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PoolID path string true "Pool ID"
// @Success 200 {object} pb.DeletePoolResponse
// @ID delete_pool
// @Router /pool-transactions/pools/{PoolID} [delete]
func DeletePool(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	poolID := ctx.Param("PoolID")

	grpcCtx := ctx.Request.Context()

	res, err := c.DeletePool(grpcCtx, &pb.DeletePoolRequest{
		PoolId: poolID,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
