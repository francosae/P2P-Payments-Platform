package routes

import (
	"fmt"
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// @Summary Creates a user record on the user-compliance service DB.
// @Description Creates a user record on the user-compliance service DB.
// @Tags User & Compliance Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} pb.CreateUserRecordResponse
// @ID create_user_record
// @Router /user-compliance/users/create-user-record/{id} [post]
func CreateUserRecord(ctx *gin.Context, c pb.UserComplianceServiceClient) {

	// TODO: Fix the param & metadata issue here. CreateUserRecord should be handle with a new account creation flow.

	userID := ctx.Param("id")

	// Create a new context with the user ID in the metadata
	md := metadata.New(map[string]string{"user_id": userID})
	grpcCtx := metadata.NewOutgoingContext(ctx.Request.Context(), md)

	res, err := c.CreateUserRecord(grpcCtx, &pb.CreateUserRecordRequest{})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, fmt.Errorf("error creating user: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

// func CreateUserRecord(ctx *gin.Context, c pb.UserComplianceServiceClient) {
// 	userID := ctx.Param("id")
// 	ctx.Set("userID", userID)

// 	grpcCtx := ctx.Request.Context()

// 	res, err := c.CreateUserRecord(grpcCtx, &pb.CreateUserRecordRequest{})

// 	if err != nil {
// 		ctx.AbortWithError(http.StatusFailedDependency, fmt.Errorf("error creating user: %v", err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, &res)
// }
