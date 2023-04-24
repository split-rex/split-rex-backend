package routes

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/configs/middlewares"
	controllers "split-rex-backend/controllers/statistic"

	"github.com/labstack/echo/v4"
)

func StatisticRoute(e *echo.Echo) {
	statisticController := controllers.NewStatisticController(database.DB.GetConnection())

	e.GET("/owedLentPercentage", statisticController.OwedLentPercentage, middlewares.AuthMiddleware)
	e.GET("/paymentMutation", statisticController.PaymentMutation, middlewares.AuthMiddleware)
}
