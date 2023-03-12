package routes

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/friend"

	"github.com/labstack/echo/v4"
)

func FriendRoute(e *echo.Echo) {
	friendController := controllers.NewFriendController(database.DB.GetConnection())

	// Get all user's friend
	e.GET("/userFriendList", friendController.UserFriendList, middlewares.AuthMiddleware)
	e.GET("/searchUser", friendController.SearchUser, middlewares.AuthMiddleware)
	e.GET("/searchUserToAdd", friendController.SearchUserToAdd, middlewares.AuthMiddleware)

	// friend request
	e.GET("/friendRequestSent", friendController.FriendRequestSent, middlewares.AuthMiddleware)
	e.GET("/friendRequestReceived", friendController.FriendRequestReceived, middlewares.AuthMiddleware)
	e.POST("/makeFriendRequest", friendController.MakeFriendRequest, middlewares.AuthMiddleware)

	// Accept and reject request received
	e.POST("/acceptRequest", friendController.AcceptRequest, middlewares.AuthMiddleware)
	e.POST("/rejectRequest", friendController.RejectRequest, middlewares.AuthMiddleware)
}
