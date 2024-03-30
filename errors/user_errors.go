package errors

import "app/pkg/apperror"

const (
	CodeUserError = 20000 + iota
	CodeUserNotFound
	CodeUserExists
	CodeUserNameExists
	CodeUserEmailExists
)

func UserNotFound() *apperror.Error {
	return apperror.NewError(CodeUserNotFound, "user not found")
}

func UserExists() *apperror.Error {
	return apperror.NewError(CodeUserExists, "user exists")
}

func UserNameExists() *apperror.Error {
	return apperror.NewError(CodeUserNameExists, "username exists")
}

func UserEmailExists() *apperror.Error {
	return apperror.NewError(CodeUserEmailExists, "email exists")
}
