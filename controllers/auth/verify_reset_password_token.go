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
	"time"

	"github.com/labstack/echo/v4"
)

func (con *authController) VerifyResetPassTokenController(c echo.Context) error {
	// check db email if exist
	db := con.db
	response := entities.Response[string]{}

	verifyTokenRequest := requests.VerifyResetPassTokenRequest{}
	if err := c.Bind(&verifyTokenRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	if verifyTokenRequest.Code == "" || verifyTokenRequest.Email == "" || verifyTokenRequest.EncryptedToken == "" {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	user := entities.User{}
	condition := entities.User{Email: verifyTokenRequest.Email}
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

	// check time difference in range 2 minutes
	timeRequested := time.Now()
	timeGenerated := userToken.TokenExpiry
	differenceInMinutes := timeRequested.Sub(timeGenerated).Minutes()
	if differenceInMinutes > 2 {
		response.Message = types.ERROR_EXPIRED_OR_INVALID_CODE
		return c.JSON(http.StatusBadRequest, response)
	}

	// check code correct
	if verifyTokenRequest.Code != userToken.Code {
		response.Message = types.ERROR_EXPIRED_OR_INVALID_CODE
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
	chiperText, err := base64.StdEncoding.DecodeString(verifyTokenRequest.EncryptedToken)
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
		response.Message = types.ERROR_EXPIRED_OR_INVALID_TOKEN
		return c.JSON(http.StatusBadRequest, response)
	}
	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
