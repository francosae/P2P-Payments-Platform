package routes

import (
	"net/http"

	"github.com/Sharefunds/api-gateway/pkg/user-compliance/pb"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UpdateUserAccountRequestBody struct {
	Username     string       `json:"username"`
	PersonalInfo PersonalInfo `json:"personalInfo"`
	Address      Address      `json:"address"`
}

// @Summary Update a user's account info on the user-compliance service DB and the corresponding user account with Galileo.
// @Description Updates a user's account info on the user-compliance service DB with provided info, currently does not support updating ID or IDType. Does not update galileo account.
// @Tags User & Compliance Service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UID path string true "User ID"
// @Param AccountInfo body UpdateUserAccountRequestBody true "Info to update user's account"
// @Success 200 {object} pb.UpdateUserAccountResponse
// @ID update_user_account
// @Router /user-compliance/users/update-account/{UID} [put]
func UpdateUserAccount(ctx *gin.Context, c pb.UserComplianceServiceClient) {
	body := UpdateUserAccountRequestBody{}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	grpcCtx := ctx.Request.Context()

	res, err := c.UpdateUserAccount(grpcCtx, &pb.UpdateUserAccountRequest{
		Username: body.Username,
		PersonalInfo: &pb.PersonalInfo{
			FirstName:   body.PersonalInfo.FirstName,
			LastName:    body.PersonalInfo.LastName,
			DateOfBirth: timestamppb.New(body.PersonalInfo.DateOfBirth),
			Id:          body.PersonalInfo.ID,
			IdType:      body.PersonalInfo.IDType,
			PhoneNumber: body.PersonalInfo.PhoneNumber,
			Email:       body.PersonalInfo.Email,
		},
		Address: &pb.Address{
			Address1:    body.Address.Address1,
			Address2:    body.Address.Address2,
			City:        body.Address.City,
			State:       body.Address.State,
			PostalCode:  body.Address.PostalCode,
			CountryCode: body.Address.CountryCode,
		},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusFailedDependency, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}
