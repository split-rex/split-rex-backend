package routes

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/transaction"

	"github.com/labstack/echo/v4"
)

func TransactionRoute(e *echo.Echo) {
	transactionController := controllers.NewTransactionController(database.DB.GetConnection())

	e.POST("/userCreateTransaction", transactionController.UserCreateTransaction, middlewares.AuthMiddleware)
	e.POST("/updatePayment", transactionController.UpdatePayment, middlewares.AuthMiddleware)
	e.POST("/resolveTransaction", transactionController.ResolveTransaction, middlewares.AuthMiddleware)
	e.GET("/getUnsettledTransaction", transactionController.GetUnsettledTransaction, middlewares.AuthMiddleware)
}
