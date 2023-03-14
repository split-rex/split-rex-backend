package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type friendController struct {
	db *gorm.DB
}

type FriendController interface {
	FriendRequestSent(c echo.Context) error
	FriendRequestReceived(c echo.Context) error
	MakeFriendRequest(c echo.Context) error
	AcceptRequest(c echo.Context) error
	RejectRequest(c echo.Context) error
	SearchUser(c echo.Context) error
	SearchUserToAdd(c echo.Context) error
	UserFriendList(c echo.Context) error
}

func NewFriendController(db *gorm.DB) FriendController {
	return &friendController{db: db}
}
