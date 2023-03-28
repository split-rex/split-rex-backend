package routes

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/payment"

	"github.com/labstack/echo/v4"
)

func PaymentRoute(e *echo.Echo) {
	paymentController := controllers.NewPaymentController(database.DB.GetConnection())

	e.POST("/updatePayment", paymentController.UpdatePayment, middlewares.AuthMiddleware)
	e.POST("/resolveTransaction", paymentController.ResolveTransaction, middlewares.AuthMiddleware)
	e.GET("/getUnsettledPayment", paymentController.GetUnsettledPayment, middlewares.AuthMiddleware)
	e.GET("/getUnconfirmedPayment", paymentController.GetUnconfirmedPayment, middlewares.AuthMiddleware)
	e.POST("/settlePaymentOwed", paymentController.SettlePaymentOwed, middlewares.AuthMiddleware)
	e.POST("/settlePaymentLent", paymentController.SettlePaymentLent, middlewares.AuthMiddleware)
	e.POST("/confirmSettle", paymentController.ConfirmSettle, middlewares.AuthMiddleware)
	e.POST("/denySettle", paymentController.DenySettle, middlewares.AuthMiddleware)
}
