package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Fetch all pools for a user.
// @Description Fetches all pools for a user.
// @Tags Pool & Transactions Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} pb.GetUserPoolsResponse
// @ID get_user_pools
// @Router /pool-transactions/pools/list [get]
func GetUserPools(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	grpcCtx := ctx.Request.Context()

	res, err := c.GetUserPools(grpcCtx, &pb.GetUserPoolsRequest{})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
