package routes

import (
	"split-rex-backend/controllers"

	"github.com/labstack/echo/v4"
)

func FriendRoute(e *echo.Echo) {
	e.POST("/makeFriendRequest", controllers.MakeFriendRequest)
}
