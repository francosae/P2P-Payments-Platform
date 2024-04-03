// @title Sharefunds Client Facing API
// @description This is the client facing API for Sharefunds/Quilt/PoolParty
// @version 1.0.1
// @host localhost:3000
// @BasePath /
// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sharefunds/api-gateway/pkg/auth"
	"github.com/Sharefunds/api-gateway/pkg/config"
	"github.com/Sharefunds/api-gateway/pkg/middleware"
	pooltransactions "github.com/Sharefunds/api-gateway/pkg/pool-transactions"
	usercompliance "github.com/Sharefunds/api-gateway/pkg/user-compliance"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed at loading config")
	}
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	r := gin.Default()
	r.Use(middleware.LogMiddleware())

	authSvc := *auth.RegisterRoutes(r, &c)
	usercompliance.RegisterRoutes(r, &c, &authSvc)
	pooltransactions.RegisterRoutes(r, &c, &authSvc)

	r.Static("/docs", "./docs")
	swaggerURL := ginSwagger.URL("/docs/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))

	fmt.Println("Sharefunds API Gateway is running")
	r.Run(c.Port)

}
