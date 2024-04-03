package pooltransactions

import (
	"github.com/Sharefunds/api-gateway/pkg/auth"
	"github.com/Sharefunds/api-gateway/pkg/config"
	"github.com/Sharefunds/api-gateway/pkg/pool-transactions/routes"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	pTMiddleware := InitPoolTransactionsMiddleware(svc)

	routes := r.Group("/pool-transactions")
	routes.Use(a.AuthRequired)
	{
		users := routes.Group("/users")
		{
			users.GET("/balance", svc.GetUserBalance)
			users.GET("/overview", svc.GetUserAccountOverview)
			users.POST("/send-p2p-payment/:UserID", svc.SendP2PPayment)
			users.POST("/send-p2p-pool-payment/:PoolID", svc.SendP2PPoolPayment)
		}

		pools := routes.Group("/pools")
		{
			pools.GET("/list", svc.GetUserPools)
			pools.GET("/balance/:PoolID", svc.GetPoolBalance)
			pools.GET("/overview/:PoolID", svc.GetPoolOverview)

			pools.GET("/:PoolID", svc.GetPoolData)
			pools.POST("/create", svc.CreatePool)
			pools.DELETE("/:PoolID", pTMiddleware.PoolOwnerRequired, svc.DeletePool)
			pools.PUT("/:PoolID", pTMiddleware.PoolOwnerRequired, svc.UpdatePool)

			pools.POST("/:PoolID/create-invite", pTMiddleware.PoolOwnerRequired, svc.CreatePoolInvitation)
			pools.POST("/:PoolID/accept-invite", svc.CreatePoolParticipant)
			pools.PUT("/:PoolID/remove-participant", pTMiddleware.PoolOwnerRequired, svc.RemovePoolParticipant)
		}
	}

}

func (svc *ServiceClient) GetPoolData(ctx *gin.Context) {
	routes.GetPoolData(ctx, svc.Client)
}

func (svc *ServiceClient) CreatePool(ctx *gin.Context) {
	routes.CreatePool(ctx, svc.Client)
}

func (svc *ServiceClient) UpdatePool(ctx *gin.Context) {
	routes.UpdatePool(ctx, svc.Client)
}

func (svc *ServiceClient) DeletePool(ctx *gin.Context) {
	routes.DeletePool(ctx, svc.Client)
}

func (svc *ServiceClient) CreatePoolInvitation(ctx *gin.Context) {
	routes.CreatePoolInvitation(ctx, svc.Client)
}

func (svc *ServiceClient) CreatePoolParticipant(ctx *gin.Context) {
	routes.CreatePoolParticipant(ctx, svc.Client)
}

func (svc *ServiceClient) RemovePoolParticipant(ctx *gin.Context) {
	routes.RemovePoolParticipant(ctx, svc.Client)
}

func (svc *ServiceClient) GetUserPools(ctx *gin.Context) {
	routes.GetUserPools(ctx, svc.Client)
}

func (svc *ServiceClient) GetPoolBalance(ctx *gin.Context) {
	routes.GetPoolBalance(ctx, svc.Client)
}

func (svc *ServiceClient) GetPoolOverview(ctx *gin.Context) {
	routes.GetPoolOverview(ctx, svc.Client)
}

func (svc *ServiceClient) SendP2PPayment(ctx *gin.Context) {
	routes.CreateP2PPayment(ctx, svc.Client)
}

func (svc *ServiceClient) SendP2PPoolPayment(ctx *gin.Context) {
	routes.CreateP2PPoolPayment(ctx, svc.Client)
}

func (svc *ServiceClient) GetUserAccountOverview(ctx *gin.Context) {
	routes.GetUserAccountOverview(ctx, svc.Client)
}

func (svc *ServiceClient) GetUserBalance(ctx *gin.Context) {
	routes.GetUserBalance(ctx, svc.Client)
}
