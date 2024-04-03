package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

type UpdatePoolRequestBody struct {
	Pool Pool `json:"Pool"`
}

// @Summary Updates pool information.
// @Description Updates pool information(name and/or description and/or goal).
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PoolID path string true "Pool ID"
// @Param Pool body UpdatePoolRequestBody true "Pool to update"
// @Success 200 {object} pb.UpdatePoolResponse
// @ID update_pool
// @Router /pool-transactions/pools/{PoolID} [put]
func UpdatePool(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	poolID := ctx.Param("PoolID")

	body := UpdatePoolRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.UpdatePool(grpcCtx, &pb.UpdatePoolRequest{
		PoolId:      poolID,
		Poolname:    body.Pool.PoolName,
		Description: body.Pool.Description,
		BalanceGoal: body.Pool.BalanceGoal,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
