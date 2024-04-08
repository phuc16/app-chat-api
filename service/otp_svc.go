package service

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type OtpService struct {
	UserRepo IUserRepo
	OtpRepo  IOtpRepo
}

func NewOtpService(userRepo IUserRepo, otpRepo IOtpRepo) *OtpService {
	return &OtpService{UserRepo: userRepo, OtpRepo: otpRepo}
}

func (s *OtpService) GenerateOtp(ctx context.Context, email string) (res entity.Otp, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	_, err = s.UserRepo.GetInactiveUser(ctx, email)
	if err != nil {
		return
	}

	otp := &entity.Otp{
		ID:        utils.NewID(),
		Email:     email,
		Code:      fmt.Sprintf("%06d", gofakeit.Number(100000, 999999)),
		CreatedAt: time.Now().Add(5 * time.Minute),
	}
	err = s.OtpRepo.SaveOtp(ctx, otp)
	if err != nil {
		return
	}
	err = s.SendOtp(ctx, email, otp.Code)
	return
}

func (s *OtpService) SendOtp(ctx context.Context, email, code string) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	message := []byte("To: " + email + "\r\n" +
		"Subject: Your OTP\r\n" +
		"\r\n" +
		"Your OTP is: " + code + "\r\n")

	auth := smtp.PlainAuth("", config.Cfg.Mail.User, config.Cfg.Mail.Password, config.Cfg.Mail.Host)
	err = smtp.SendMail(fmt.Sprintf("%s:%d", config.Cfg.Mail.Host, config.Cfg.Mail.Port), auth, config.Cfg.Mail.User, []string{email}, message)
	if err != nil {
		return err
	}
	return
}

func (s *OtpService) VerifyOtp(ctx context.Context, e *entity.Otp) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	otp, err := s.OtpRepo.GetOtp(ctx, e)
	if err != nil {
		return
	}
	if time.Now().After(otp.CreatedAt) {
		err2 := s.OtpRepo.DeleteOtp(ctx, otp)
		if err != nil {
			return res, err2
		}
		return res, errors.OtpExpired()
	}

	dbUser, err := s.UserRepo.GetInactiveUser(ctx, e.Email)
	if err != nil {
		return
	}
	dbUser.OnUserActive(ctx)
	err = s.UserRepo.UpdateUser(ctx, dbUser)
	if err != nil {
		return
	}

	err = s.OtpRepo.DeleteOtp(ctx, otp)
	return
}