package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/pb"
	"github.com/gin-gonic/gin"
)

type CreateP2PPaymentRequestBody struct {
	Amount      float32 `json:"Amount"`
	Description string  `json:"Description"`
}

// @Summary Create P2P Payment.
// @Description Sends a P2P payment from one user to another, an internal account transfer via Galileo.
// @Tags Pool & Transactions Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param PaymentInfo body CreateP2PPaymentRequestBody true "Payment information"
// @Param ReceiverID path string true "Receiver ID for payment"
// @Success 200 {object} pb.SendPaymentToUserResponse
// @ID create_p2p_payment
// @Router /pool-transactions/users/send-p2p-payment/{ReceiverID} [post]
func CreateP2PPayment(ctx *gin.Context, c pb.PoolTransactionsServiceClient) {
	body := CreateP2PPaymentRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.SendPaymentToUser(grpcCtx, &pb.SendPaymentToUserRequest{
		ReceiverUserId: ctx.Param("ReceiverID"),
		Amount:         body.Amount,
		Description:    body.Description,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
