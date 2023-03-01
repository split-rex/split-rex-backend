package controllers

import (
	"net/http"
	"split-rex-backend/configs"
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterController(c echo.Context) error {
	db := database.DB.GetConnection()
	config := configs.Config.GetMetadata()
	response := entities.Response[string]{}

	registerRequest := requests.RegisterRequest{}
	if err := c.Bind(&registerRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	user := entities.User{}
	conditionEmail := entities.User{Email: registerRequest.Email}
	conditionUsername := entities.User{Username: registerRequest.Username}

	// check email
	if err := db.Where(&conditionEmail).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	if user.Email != "" {
		response.Message = types.ERROR_EMAIL_EXISTED
		return c.JSON(http.StatusBadRequest, response)
	}

	// check username
	if err := db.Where(&conditionUsername).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	if user.Username != "" {
		response.Message = types.ERROR_USERNAME_EXISTED
		return c.JSON(http.StatusBadRequest, response)
	}

	// insert user
	if err := db.Create(&entities.User{
		ID:       uuid.New(),
		Username: registerRequest.Username,
		Email:    registerRequest.Email,
		Name:     registerRequest.Name,
		Password: types.EncryptedString(registerRequest.Password),
	}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	unsignedAuthToken := jwt.NewWithClaims(config.JWTSigningMethod, middlewares.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginExpirationDuration)),
		},
		ID: user.ID,
	})

	signedAuthToken, err := unsignedAuthToken.SignedString(config.JWTSignatureKey)
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	response.Data = signedAuthToken
	return c.JSON(http.StatusAccepted, response)
}
