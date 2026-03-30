package errors

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Msg string `json:"message"`
	Code int `json:"code"`
	Err error `json:"-"`
}

func (e HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Msg, e.Err)
	}
	return e.Msg
}

func BadRequest(msg string) HTTPError {
	return HTTPError{
		Msg: msg,
		Code: http.StatusBadRequest,
	}
}

func Unauthorized(msg string) HTTPError {
	return HTTPError{
		Msg: msg,
		Code: http.StatusUnauthorized,
	}
}

func InternalServerError(err error) HTTPError {
	return HTTPError{
		Msg: "internal server error",
		Code: http.StatusInternalServerError,
		Err: err,
	}
}
