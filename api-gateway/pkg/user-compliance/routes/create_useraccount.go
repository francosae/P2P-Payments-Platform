package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
)

type CreateUserAccountRequestBody struct {
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Address      Address      `json:"address"`
	Username     string       `json:"username"`
}

type PersonalInfo struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber string    `json:"phoneNumber"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"dateOfBirth,omitempty"`
	ID          string    `json:"id"`
	IDType      int32     `json:"idType"`
}

type Address struct {
	Address1    string `json:"address1"`
	Address2    string `json:"address2,omitempty"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postalCode"`
	CountryCode string `json:"countryCode"`
}

// @Summary Create a user's digital wallet (galileo) account, and a record of the users info on the user-compliance service DB.
// @Description Creates a user's Galileo account, and account record within our user-compliance service DB with provided info, returns the Galileo account info (Create a new card program) and status.
// @Tags User & Compliance Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UID path string true "User ID"
// @Param AccountInfo body CreateUserAccountRequestBody true "Request Body to create user account"
// @Success 200 {object} pb.CreateUserAccountResponse
// @ID create_user_account
// @Router /user-compliance/users/create-account/{UID} [post]
func CreateUserAccount(ctx *gin.Context, c pb.UserComplianceServiceClient) {
	body := CreateUserAccountRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.CreateUserAccount(grpcCtx, &pb.CreateUserAccountRequest{
		Username: body.Username,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, fmt.Errorf("error creating user: %v", err))
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
