package routes

import (
	"split-rex-backend/configs"
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/auth"

	"github.com/labstack/echo/v4"
)

func AuthRoute(e *echo.Echo) {
	authController := controllers.NewAuthController(database.DB.GetConnection(), configs.Config.GetMetadata())

	e.POST("/login", authController.LoginController)
	e.POST("/register", authController.RegisterController)
	e.GET("/profile", authController.ProfileController, middlewares.AuthMiddleware)
}
