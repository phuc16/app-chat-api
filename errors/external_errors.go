package errors

import "app/pkg/apperror"

const (
	CodeExternalError = 80000 + iota
)

func ExternalError() *apperror.Error {
	return apperror.NewError(CodeExternalError, "external error")
}

func WrapExternalError(err *error) {
	if *err != nil {
		if apperror.As(*err) != nil {
			return
		}
		*err = apperror.NewError(CodeExternalError, (*err).Error())
	}
}
