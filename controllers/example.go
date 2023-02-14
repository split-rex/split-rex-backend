package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ExampleController(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello")
}
