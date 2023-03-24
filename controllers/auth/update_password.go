package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (con *authController) UpdatePasswordController(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	updatePasswordRequest := requests.UpdatePasswordRequest{}
	if err := c.Bind(&updatePasswordRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	// get id from context, cast to uuid
	user_id := c.Get("id").(uuid.UUID)

	// search the id of the user
	user := entities.User{}
	conditionUser := entities.User{ID: user_id}
	if err := db.Where(&conditionUser).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	// error if user not found
	if user.ID == uuid.Nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}
	// check the password of user
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(updatePasswordRequest.OldPassword)); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// update user
	if err := db.Model(&user).Updates(entities.User{
		Password: types.EncryptedString(updatePasswordRequest.NewPassword),
	}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
