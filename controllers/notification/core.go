package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type notificationController struct {
	db *gorm.DB
}

type NotificationController interface {
	GetNotif(c echo.Context) error
	InsertNotif(c echo.Context) error
	DeleteNotif(c echo.Context) error
}

func NewNotificationController(db *gorm.DB) NotificationController {
	return &notificationController{db: db}
}
