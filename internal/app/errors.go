package app

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidToken   = errors.New("invalid token")
	ErrInvalidUserId  = errors.New("invalid user id")
)

type Error struct {
	code int
	err  error
}

func NewError(code int, err error) *Error {
	return &Error{
		code: code,
		err:  err,
	}
}

func NewErrorf(code int, format string, args ...interface{}) *Error {
	return &Error{
		code: code,
		err:  fmt.Errorf(format, args...),
	}
}

func NewInternalServerError(err error) *Error {
	return NewError(http.StatusInternalServerError, err)
}

func NewBadRequestError(err error) *Error {
	return NewError(http.StatusBadRequest, err)
}

func (e Error) Code() int {
	if e.code == 0 {
		return http.StatusInternalServerError
	}

	return e.code
}

func (e Error) Error() string {
	return e.err.Error()
}
