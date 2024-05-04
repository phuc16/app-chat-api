package errors

import "app/pkg/apperror"

const (
	CodeAuthError = 80000 + iota
	CodePermissionDenied
)

func PermissionDenied() *apperror.Error {
	return apperror.NewError(CodePermissionDenied, "permission denied")
}
