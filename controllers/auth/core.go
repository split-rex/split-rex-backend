package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type authController struct {
	db *gorm.DB
}

type AuthController interface {
	LoginController(c echo.Context) error
	ProfileController(c echo.Context) error
	RegisterController(c echo.Context) error
}

func NewAuthController(db *gorm.DB) AuthController {
	return &authController{db: db}
}
