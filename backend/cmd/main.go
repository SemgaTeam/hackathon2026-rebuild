package main

import (
	e "github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	echo := e.New()

	echo.GET("/", func (c e.Context) error {
		return c.String(http.StatusOK, "hello world")
	})

	echo.Start("0.0.0.0")
}
