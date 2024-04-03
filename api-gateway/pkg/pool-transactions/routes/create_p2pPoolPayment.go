package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

type CreateP2PPoolPaymentRequestBody struct {
	Amount      float32 `json:"Amount"`
	Description string  `json:"Description"`
}

// @Summary Create P2P Payment to a pool.
// @Description Sends a P2P payment from one user to a pool, an internal account transfer via Galileo.
// @Tags Pool & Transactions Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PaymentInfo body CreateP2PPoolPaymentRequestBody true "Payment information"
// @Param PoolID path string true "Pool ID for payment"
// @Success 200 {object} pb.SendPaymentToUserResponse
// @ID create_p2p_pool_payment
// @Router /pool-transactions/users/send-p2p-pool-payment/{PoolID} [post]
func CreateP2PPoolPayment(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	body := CreateP2PPaymentRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.SendPaymentToPool(grpcCtx, &pb.SendPaymentToPoolRequest{
		PoolId:      ctx.Param("PoolID"),
		Amount:      body.Amount,
		Description: body.Description,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
