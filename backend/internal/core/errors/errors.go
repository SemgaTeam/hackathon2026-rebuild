package errors

import (
	"fmt"
)

type DomainError struct {
	Msg string
	Err error
}

func NewError(msg string) DomainError {
	return DomainError{
		Msg: msg,
	}
}

func Unknown(err error) DomainError {
	return DomainError{
		Msg: "unknown error",
		Err: err,
	}
}

func (e DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e DomainError) Unwrap() error {
	return e.Err
}

var (
	// ErrUserNotFound = NewError("user not found")
)
