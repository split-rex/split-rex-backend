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

	e.POST("/generateResetPassToken", authController.GenerateResetPassTokenController);
	e.POST("/verifyResetPassToken", authController.VerifyResetPassTokenController);
	e.POST("/changePassword", authController.ChangePasswordController);

	e.POST("/login", authController.LoginController)
	e.POST("/register", authController.RegisterController)

	e.GET("/profile", authController.ProfileController, middlewares.AuthMiddleware)
	e.POST("/updateProfile", authController.UpdateProfileController, middlewares.AuthMiddleware)
	e.POST("/updatePassword", authController.UpdatePasswordController, middlewares.AuthMiddleware)

	e.POST("/addPaymentInfo", authController.AddPaymentInfo, middlewares.AuthMiddleware)
	e.POST("/editPaymentInfo", authController.EditPaymentInfo, middlewares.AuthMiddleware)
	e.POST("/deletePaymentInfo", authController.DeletePaymentInfo, middlewares.AuthMiddleware)
}
