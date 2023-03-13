package controllers

import (
	"split-rex-backend/configs"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type authController struct {
	db       *gorm.DB
	metadata configs.Metadata
}

type AuthController interface {
	LoginController(c echo.Context) error
	ProfileController(c echo.Context) error
	RegisterController(c echo.Context) error
}

func NewAuthController(db *gorm.DB, mt configs.Metadata) AuthController {
	return &authController{db: db, metadata: mt}
}
