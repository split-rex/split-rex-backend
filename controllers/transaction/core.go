package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type transactionController struct {
	db *gorm.DB
}

type TransactionController interface {
	UserCreateTransaction(c echo.Context) error
}

func NewTransactionController(db *gorm.DB) TransactionController {
	return &transactionController{db: db}
}
