package routes

import (
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/transaction"

	"github.com/labstack/echo/v4"
)

func TransactionRoute(e *echo.Echo) {
	e.POST("/userCreateTransaction", controllers.UserCreateTransactionController, middlewares.AuthMiddleware)
}
