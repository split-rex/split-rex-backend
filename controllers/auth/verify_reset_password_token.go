package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"split-rex-backend/configs"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"
	"time"

	"github.com/labstack/echo/v4"
)

func (con *authController) VerifyResetPassTokenController(c echo.Context) error {
	// check db email if exist
	config := configs.Config.GetMetadata()

	db := con.db
	response := entities.Response[responses.VerifyResetPassTokenResponse]{}

	verifyTokenRequest := requests.VerifyResetPassTokenRequest{}
	if err := c.Bind(&verifyTokenRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	if verifyTokenRequest.Code == "" || verifyTokenRequest.Email == ""{
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

	// generate random for passwordToken
	key := config.ResetPasswordKey

	// generate a new aes cipher using our 32 byte long key
	cip, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	gcm, err := cipher.NewGCM(cip)
	// if any error generating new GCM
	// handle them
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	token := userToken.Token
	encryptedToken := base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(token), nil))

	response.Data.EncryptedToken = encryptedToken
	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
