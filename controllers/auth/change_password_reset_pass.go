package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/http"
	"os"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/labstack/echo/v4"
)

func (con *authController) ChangePasswordController(c echo.Context) error {

	// check db if exist
	db := con.db
	response := entities.Response[string]{}

	changePasswordRequest := requests.ChangePasswordRequest{}
	if err := c.Bind(&changePasswordRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	if changePasswordRequest.Code == "" || changePasswordRequest.Email == "" || changePasswordRequest.EncryptedToken == "" || changePasswordRequest.NewPassword == "" {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	user := entities.User{}
	condition := entities.User{Email: changePasswordRequest.Email}
	if err := db.Where(&condition).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if user.Email == "" {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	userToken := entities.PasswordResetTokens{}
	condition2 := entities.PasswordResetTokens{Email: user.Email}
	if err := db.Where(&condition2).Find(&userToken).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if userToken.Token == "" {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}
	// check code correct
	if changePasswordRequest.Code != userToken.Code {
		response.Message = types.ERROR_INVALID_CODE
		return c.JSON(http.StatusBadRequest, response)
	}

	// check token correct
	key := []byte(os.Getenv("RESET_PASS_KEY"))
	cip, err := aes.NewCipher(key)
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	gcm, err := cipher.NewGCM(cip)
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	chiperText, err := base64.StdEncoding.DecodeString(changePasswordRequest.EncryptedToken)
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	nonceSize := gcm.NonceSize()
	if len(chiperText) < nonceSize {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	nonce, chiperText := chiperText[:nonceSize], chiperText[nonceSize:]
	decryptedToken, err := gcm.Open(nil, nonce, chiperText, nil)
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if string(decryptedToken) != userToken.Token {
		response.Message = types.ERROR_INVALID_TOKEN
		return c.JSON(http.StatusBadRequest, response)
	}
	// change the password in user db
	if err := db.Model(&user).Updates(entities.User{
		Password: types.EncryptedString(changePasswordRequest.NewPassword),
	}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// delete the row in password_reset_token
	if err := db.Where(&entities.PasswordResetTokens{
		Email: userToken.Email,
		Code:  userToken.Code,
		Token: userToken.Token}).Delete(&entities.PasswordResetTokens{}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
