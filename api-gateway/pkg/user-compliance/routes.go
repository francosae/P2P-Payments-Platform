package usercompliance

import (
	"github.com/Sharefunds/api-gateway/pkg/auth"
	"github.com/Sharefunds/api-gateway/pkg/config"
	"github.com/Sharefunds/api-gateway/pkg/user-compliance/routes"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	routes := r.Group("/user-compliance")

	// TODO: Change record creation to be done when a user is created through Firebase on server side.
	routes.POST("/users/create-user-record/:id", svc.CreateUserRecord)

	routes.Use(a.AuthRequired)
	{
		users := routes.Group("/users")
		users.Use(a.UserIdRequired())
		{
			users.GET("/get-account/:id", svc.GetUserAccount)
			users.PUT("/update-account/:id", svc.UpdateUserAccount)
			users.POST("/create-account/:id", svc.CreateUserAccount)
			users.POST("/set-verification-status/:id", svc.SetIdentityStatus)
			users.GET("/get-account-status/:id", svc.GetAccountStatus)
		}

		plaid := routes.Group("/plaid")
		{
			plaid.POST("/create-linktoken/:id", svc.CreateLinkToken)
			plaid.POST("/create-accesstoken/:id", svc.CreateAccessToken)
			plaid.POST("/create-idvtoken/:id", svc.CreateIDVToken)
		}
	}

}

func (svc *ServiceClient) CreateUserAccount(ctx *gin.Context) {
	routes.CreateUserAccount(ctx, svc.Client)

}

func (svc *ServiceClient) UpdateUserAccount(ctx *gin.Context) {
	routes.UpdateUserAccount(ctx, svc.Client)
}

func (svc *ServiceClient) GetUserAccount(ctx *gin.Context) {
	routes.GetUserAccount(ctx, svc.Client)
}

// Plaid
func (svc *ServiceClient) CreateLinkToken(ctx *gin.Context) {
	routes.CreateLinkToken(ctx, svc.Client)
}

func (svc *ServiceClient) CreateAccessToken(ctx *gin.Context) {
	routes.CreateAccessToken(ctx, svc.Client)
}

func (svc *ServiceClient) CreateIDVToken(ctx *gin.Context) {
	routes.CreateIDVToken(ctx, svc.Client)
}

// Identity Verification
func (svc *ServiceClient) SetIdentityStatus(ctx *gin.Context) {
	routes.SetIdentityStatus(ctx, svc.Client)
}

func (svc *ServiceClient) GetAccountStatus(ctx *gin.Context) {
	routes.GetAccountStatus(ctx, svc.Client)
}

func (svc *ServiceClient) CreateUserRecord(ctx *gin.Context) {
	routes.CreateUserRecord(ctx, svc.Client)
}
