package routes

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/payment"

	"github.com/labstack/echo/v4"
)

func PaymentRoute(e *echo.Echo) {
	transactionController := controllers.NewPaymentController(database.DB.GetConnection())

	e.POST("/updatePayment", transactionController.UpdatePayment, middlewares.AuthMiddleware)
	e.POST("/resolveTransaction", transactionController.ResolveTransaction, middlewares.AuthMiddleware)
	e.GET("/getUnsettledPayment", transactionController.GetUnsettledPayment, middlewares.AuthMiddleware)
	e.GET("/getUnconfirmedPayment", transactionController.GetUnconfirmedPayment, middlewares.AuthMiddleware)
	e.POST("/settlePaymentOwed", transactionController.SettlePaymentOwed, middlewares.AuthMiddleware)
	e.POST("/settlePaymentLent", transactionController.SettlePaymentLent, middlewares.AuthMiddleware)
}
