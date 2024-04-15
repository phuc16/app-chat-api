package errors

import "app/pkg/apperror"

const (
	CodeUserError = 20000 + iota
	CodeUserNotFound
	CodeUserExists
	CodeUserNameExists
	CodeUserEmailExists
	CodePasswordIncorrect
	CodeUserNotRegister
	CodeUserInactive
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

func PasswordIncorrect() *apperror.Error {
	return apperror.NewError(CodePasswordIncorrect, "password is incorrect")
}

func UserNotRegister() *apperror.Error {
	return apperror.NewError(CodeUserNotRegister, "You do not have an account, please register first")
}

func UserInactive() *apperror.Error {
	return apperror.NewError(CodeUserInactive, "account has not been active, please verify your email first")
}
