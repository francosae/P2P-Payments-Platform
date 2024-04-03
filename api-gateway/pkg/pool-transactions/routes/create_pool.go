package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

type Pool struct {
	PoolName    string `json:"poolName"`
	Description string `json:"description"`
	BalanceGoal int32  `json:"balanceGoal"`
}

type CreatePoolRequestBody struct {
	Pool Pool `json:"Pool"`
}

// @Summary Create a new pool
// @Description Creates a new pool in the database, and creates a new account in Galileo for the pool, and adds the user as a participant in the pool.
// @Tags Pool & Transactions Service - Pools
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Pool body CreatePoolRequestBody true "Pool to create"
// @Success 200 {object} pb.CreatePoolResponse
// @ID create_pool
// @Router /pool-transactions/pools/create [post]
func CreatePool(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	body := CreatePoolRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.CreatePool(grpcCtx, &pb.CreatePoolRequest{
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
