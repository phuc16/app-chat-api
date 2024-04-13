package dto

import (
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"context"

	"github.com/gin-gonic/gin"
)

type OtpReq struct {
	Email string `json:"email" binding:"required"`
}

func (r OtpReq) Bind(ctx *gin.Context) (res *OtpReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}
func (r OtpReq) Validate() (err error) {
	return
}

func (r OtpReq) ToOtp(ctx context.Context) (res *entity.Otp) {
	res = &entity.Otp{
		Email: r.Email,
	}
	return res
}
