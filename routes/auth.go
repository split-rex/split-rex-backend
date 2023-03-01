package routes

import (
	controllers "split-rex-backend/controllers/auth"

	"github.com/labstack/echo/v4"
)

func AuthRoute(e *echo.Echo) {
	e.POST("/login", controllers.LoginController)
	e.POST("/register", controllers.RegisterController)

	e.GET("/profile", controllers.ProfileController)
}
