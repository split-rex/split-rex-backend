package routes

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/activity"

	"github.com/labstack/echo/v4"
)

func ActivityRoute(e *echo.Echo) {
	activityController := controllers.NewActivityController(database.DB.GetConnection())

	e.GET("/getUserActivity", activityController.GetUserActivity, middlewares.AuthMiddleware)
	e.GET("/getGroupActivity", activityController.GetGroupActivity, middlewares.AuthMiddleware)
}
