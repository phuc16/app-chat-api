package errors

import "app/pkg/apperror"

const (
	CodeUnknownError = 10000 + iota
)

func WrapError(err *error) {
	if *err != nil {
		if apperror.As(*err) != nil {
			return
		}
		*err = apperror.NewError(CodeUnknownError, (*err).Error())
	}
}
