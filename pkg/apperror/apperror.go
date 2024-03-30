package apperror

import (
	"fmt"
	"github.com/pkg/errors"
)

type Error struct {
	Code int
	Err  error
}

func NewError(code int, msg string) *Error {
	return &Error{Code: code, Err: errors.New(msg)}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Wrap(msg string) *Error {
	e.Err = errors.Wrap(e.Err, msg)
	return e
}

func (e *Error) IsNot(err error) bool {
	return !e.Is(err)
}

func (e *Error) Is(err error) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		if err.(*Error).Code == e.Code {
			return true
		}
	}
	return false
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (e *Error) StackTrace() string {
	err, ok := errors.Cause(e.Err).(stackTracer)
	if !ok {
		return e.Error()
	}
	st := err.StackTrace()
	return fmt.Sprintf("%+v\n", st[1:5])
}

func As(err error) *Error {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
