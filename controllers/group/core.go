package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type groupController struct {
	db *gorm.DB
}

type GroupController interface {
	UserCreateGroup(c echo.Context) error
	EditGroupInfo(c echo.Context) error
	UserGroups(c echo.Context) error
	GroupDetail(c echo.Context) error
	GroupTransactions(c echo.Context) error
	GroupLent(c echo.Context) error
	GroupOwed(c echo.Context) error
}

func NewGroupController(db *gorm.DB) GroupController {
	return &groupController{db: db}
}
