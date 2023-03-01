package routes

import (
	"split-rex-backend/configs/middlewares"
	"split-rex-backend/controllers"

	"github.com/labstack/echo/v4"
)

func FriendRoute(e *echo.Echo) {
	// Get all user's friend
	e.GET("/userFriendList", controllers.UserFriendList, middlewares.AuthMiddleware)
	e.GET("/friendRequestSent", controllers.FriendRequestSent, middlewares.AuthMiddleware)
	e.GET("/friendRequestReceived", controllers.FriendRequestReceived, middlewares.AuthMiddleware)
	e.GET("/searchUser", controllers.SearchUser, middlewares.AuthMiddleware)

	e.POST("/makeFriendRequest", controllers.MakeFriendRequest, middlewares.AuthMiddleware)

	// Accept and reject request received
	e.POST("/acceptRequest", controllers.AcceptRequest, middlewares.AuthMiddleware)
	e.POST("/rejectRequest", controllers.RejectRequest, middlewares.AuthMiddleware)
}
