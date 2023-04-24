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
}

func NewStatisticController(db *gorm.DB) StatisticController {
	return &statisticController{db: db}
}
