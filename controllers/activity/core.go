package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type activityController struct {
	db *gorm.DB
}

type ActivityController interface {
	GetUserActivity(c echo.Context) error
	GetGroupActivity(c echo.Context) error
}

func NewActivityController(db *gorm.DB) ActivityController {
	return &activityController{db: db}
}
