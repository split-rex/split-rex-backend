package routes

import (
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/auth"

	"github.com/labstack/echo/v4"
)

func TransactionRoute(e *echo.Echo) {
	e.POST("/userCreateTransaction", controllers.LoginController, middlewares.AuthMiddleware)
}
