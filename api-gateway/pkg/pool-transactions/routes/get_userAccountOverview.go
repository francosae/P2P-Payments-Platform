package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

// @Summary Get User Account Overview.
// @Description Fetches the overview of a user's account, including user's balance, transactions, and transaction count.
// @Tags Pool & Transactions Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UserID path string true "User ID"
// @Success 200 {object} pb.GetUserAccountOverviewResponse
// @ID get_user_account_overview
// @Router /pool-transactions/users/overview/{UserID} [get]
func GetUserAccountOverview(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	grpcCtx := ctx.Request.Context()

	res, err := c.GetUserAccountOverview(grpcCtx, &pb.GetUserAccountOverviewRequest{
		UserId: ctx.Param("UserID"),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
