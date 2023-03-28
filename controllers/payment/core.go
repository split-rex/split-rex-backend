package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type paymentController struct {
	db *gorm.DB
}

type PaymentController interface {
	UpdatePayment(c echo.Context) error
	GetUnsettledTransaction(c echo.Context) error
	ResolveTransaction(c echo.Context) error
}

func NewPaymentController(db *gorm.DB) PaymentController {
	return &paymentController{db: db}
}
