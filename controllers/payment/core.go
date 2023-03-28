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
	GetUnsettledPayment(c echo.Context) error
	ResolveTransaction(c echo.Context) error
	GetUnconfirmedPayment(c echo.Context) error
	SettlePaymentOwed(c echo.Context) error
	SettlePaymentLent(c echo.Context) error
}

func NewPaymentController(db *gorm.DB) PaymentController {
	return &paymentController{db: db}
}
