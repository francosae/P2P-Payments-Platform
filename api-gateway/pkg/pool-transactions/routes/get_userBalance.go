package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Get's user balance.
// @Description Fetches the balance of a user, using the userId.
// @Tags Pool & Transactions Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UserID path string true "User ID"
// @Success 200 {object} pb.GetUserBalanceResponse
// @ID get_user_balance
// @Router /pool-transactions/users/balance/{UserID} [get]
func GetUserBalance(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	grpcCtx := ctx.Request.Context()

	res, err := c.GetUserBalance(grpcCtx, &pb.GetUserBalanceRequest{
		UserId: ctx.Param("UserID"),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
