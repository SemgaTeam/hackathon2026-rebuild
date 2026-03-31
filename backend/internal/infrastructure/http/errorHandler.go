package http

import (
	"errors"

	e "github.com/SemgaTeam/semga-stream/internal/infrastructure/http/errors"
	"github.com/labstack/echo/v4"
)

func errorHandler(err error, c echo.Context) {
	var httpErr e.HTTPError
	var echoErr *echo.HTTPError

	switch {
	case errors.As(err, &echoErr):
		httpErr.Code = echoErr.Code	
		httpErr.Msg = echoErr.Message.(string)
	default:
		httpErr = e.InternalServerError(err)
	}

	if !c.Response().Committed {
		c.JSON(httpErr.Code, map[string]string{
			"error": httpErr.Msg,
		})
	}
}
