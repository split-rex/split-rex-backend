package core

import (
	"os"
	"split-rex-backend/configs/middlewares"
	"split-rex-backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {
	e := echo.New()

	/* Middlewares */
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.Cors())

	/* Routes */
	routes.AuthRoute(e)
	routes.GroupRoute(e)
	routes.TransactionRoute(e)
	routes.PaymentRoute(e)
	routes.FriendRoute(e)
	routes.GroupRoute(e)
	routes.ActivityRoute(e)
	routes.StatisticRoute(e)
	routes.NotificationRoute(e)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
