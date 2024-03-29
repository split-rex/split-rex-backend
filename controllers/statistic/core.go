package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type statisticController struct {
	db *gorm.DB
}

type StatisticController interface {
	OwedLentPercentage(c echo.Context) error
	PaymentMutation(c echo.Context) error
	SpendingBuddies(c echo.Context) error
	ExpenseChart(c echo.Context) error
}

func NewStatisticController(db *gorm.DB) StatisticController {
	return &statisticController{db: db}
}
