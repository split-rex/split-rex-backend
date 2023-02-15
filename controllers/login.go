package controllers

import (
	"net/http"
	"split-rex-backend/configs"
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
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
		response.Message = "ERROR: BAD REQUEST"
		return c.JSON(http.StatusBadRequest, response)
	}

	user := entities.User{}
	condition := entities.User{Email: loginRequest.Email}
	if err := db.Where(&condition).Find(&user).Error; err != nil {
		response.Message = "ERROR: INTERNAL SERVER ERROR"
		return c.JSON(http.StatusInternalServerError, response)
	}
	if user.Username == "" {
		response.Message = "ERROR: INVALID USERNAME OR PASSWORD"
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(loginRequest.Password)); err != nil {
		response.Message = "ERROR: INVALID USERNAME OR PASSWORD"
		return c.JSON(http.StatusUnauthorized, response)
	}

	unsignedAuthToken := jwt.NewWithClaims(config.JWTSigningMethod, middlewares.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.LoginExpirationDuration)),
		},
		ID:   user.ID,
	})
	
	signedAuthToken, err := unsignedAuthToken.SignedString(config.JWTSignatureKey)
	if err != nil {
		response.Message = "ERROR: JWT SIGNING ERROR"
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = "SUCCESS"
	response.Data = signedAuthToken
	return c.JSON(http.StatusAccepted, response)
}
