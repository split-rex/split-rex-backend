package routes

import (
	"split-rex-backend/controllers"

	"github.com/labstack/echo/v4"
)

func ExampleRoute(e *echo.Echo) {
	e.POST("/", controllers.ExampleController)
}
