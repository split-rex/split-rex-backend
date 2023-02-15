package routes

import (
	"split-rex-backend/controllers"

	"github.com/labstack/echo/v4"
)

func AuthRoute(e *echo.Echo) {
	e.POST("/login", controllers.LoginController)
	e.POST("/register", controllers.RegisterController)
}
