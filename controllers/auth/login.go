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
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func LoginController(c echo.Context) error {
	db := database.DB.GetConnection()
	config := configs.Config.GetMetadata()
	response := entities.Response[string]{}

	loginRequest := requests.LoginRequest{}
	if err := c.Bind(&loginRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	user := entities.User{}
	condition := entities.User{Email: loginRequest.Email}
	if err := db.Where(&condition).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	if user.Username == "" {
		response.Message = types.ERROR_FAILED_LOGIN
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(loginRequest.Password)); err != nil {
		response.Message = types.ERROR_FAILED_LOGIN
		return c.JSON(http.StatusUnauthorized, response)
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
		response.Message = types.ERROR_JWT_SIGNING
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	response.Data = signedAuthToken
	return c.JSON(http.StatusAccepted, response)
}

func LoginGoogleHandler(c echo.Context) {

}
