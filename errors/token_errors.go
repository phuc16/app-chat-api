package errors

import "app/pkg/apperror"

const (
	CodeTokenError = 30000 + iota
	CodeTokenNotFound
	CodeTokenExists
)

func TokenNotFound() *apperror.Error {
	return apperror.NewError(CodeTokenNotFound, "token not found")
}

func TokenExists() *apperror.Error {
	return apperror.NewError(CodeTokenExists, "token exists")
}
