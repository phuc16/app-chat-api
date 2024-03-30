package errors

import "app/pkg/apperror"

const (
	CodeDatabaseError = 90000 + iota
)

func DatabaseError() *apperror.Error {
	return apperror.NewError(CodeDatabaseError, "database error")
}

func WrapDatabaseError(err *error) {
	if *err != nil {
		if apperror.As(*err) != nil {
			return
		}
		*err = apperror.NewError(CodeDatabaseError, (*err).Error())
	}
}
