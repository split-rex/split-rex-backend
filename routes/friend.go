package routes

import (
	"split-rex-backend/configs/middlewares"
	"split-rex-backend/controllers"

	"github.com/labstack/echo/v4"
)

func FriendRoute(e *echo.Echo) {
	e.POST("/makeFriendRequest", controllers.MakeFriendRequest, middlewares.AuthMiddleware)
	e.GET("/friendRequestSent", controllers.FriendRequestSent, middlewares.AuthMiddleware)
	e.GET("/friendRequestReceived", controllers.FriendRequestReceived, middlewares.AuthMiddleware)
	e.GET("/searchUser", controllers.SearchUser, middlewares.AuthMiddleware)
}
