package routes

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/notification"

	"github.com/labstack/echo/v4"
)

func NotificationRoute(e *echo.Echo) {
	notificationController := controllers.NewNotificationController(database.DB.GetConnection())

	e.GET("/getNotif", notificationController.GetNotif, middlewares.AuthMiddleware)
	e.POST("/insertNotif", notificationController.InsertNotif, middlewares.AuthMiddleware)
	e.POST("/deleteNotif", notificationController.DeleteNotif, middlewares.AuthMiddleware)
}
