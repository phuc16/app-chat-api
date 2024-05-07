package errors

import "app/pkg/apperror"

const (
	CodeOtpError = 40000 + iota
	CodeOtpNotFound
	CodeOtpExists
	CodeOtpExpired
)

func OtpNotFound() *apperror.Error {
	return apperror.NewError(CodeOtpNotFound, "otp not found")
}

func OtpExists() *apperror.Error {
	return apperror.NewError(CodeOtpExists, "otp exists")
}

func OtpExpired() *apperror.Error {
	return apperror.NewError(CodeOtpExpired, "otp expired")
}
