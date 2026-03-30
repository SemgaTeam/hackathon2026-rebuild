package main

import (
	"github.com/SemgaTeam/semga-stream/internal/infrastructure/http"
	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/labstack/echo/v4"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	http.SetupHandlers(conf, e)
}
