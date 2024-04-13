package service

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"app/pkg/trace"
	"app/pkg/utils"
	"app/repository"
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserService struct {
	OtpSvc    IOtpSvc
	UserRepo  IUserRepo
	TokenRepo ITokenRepo
}

func NewUserService(otpSvc IOtpSvc, userRepo IUserRepo, tokenRepo ITokenRepo) *UserService {
	return &UserService{OtpSvc: otpSvc, UserRepo: userRepo, TokenRepo: tokenRepo}
}

func (s *UserService) Login(ctx context.Context, user *entity.User) (accessToken string, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	dbUser, err := s.UserRepo.GetUserByUserNameOrEmail(ctx, user.Username, user.Email)
	if err != nil {
		return
	}
	if !dbUser.IsActive {
		return "", errors.UserInactive()
	}
	if !utils.VerifyPassword(user.Password, dbUser.Password) {
		err = errors.PasswordIncorrect()
		return
	}
	accessToken, err = s.CreateToken(ctx, dbUser)
	if err != nil {
		return
	}
	dbUser.LoggedIn()
	_ = s.UserRepo.SaveUser(ctx, dbUser)
	return
}

func (s *UserService) CreateToken(ctx context.Context, user *entity.User) (accessToken string, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	duration := time.Duration(int32(config.Cfg.HTTP.AccessTokenDuration)) * time.Minute
	appToken := &entity.Token{
		ID:       utils.NewID(),
		UserID:   user.ID,
		Name:     user.Name,
		UserName: user.Username,
		Email:    user.Email,
		Type:     entity.AccessTokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	err = s.TokenRepo.CreateToken(ctx, appToken)
	if err != nil {
		return
	}
	accessToken = appToken.SignedToken(config.Cfg.HTTP.Secret)
	return
}

func (s *UserService) UserLogout(ctx context.Context, tokenStr string) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	token, err := entity.NewTokenFromEncoded(tokenStr, config.Cfg.HTTP.Secret)
	if err != nil {
		err = apperror.NewError(errors.CodeTokenError, err.Error())
		return
	}
	dbToken, err := s.TokenRepo.GetTokenById(ctx, token.ID)
	if err != nil {
		return
	}
	err = s.TokenRepo.DeleteToken(ctx, dbToken)
	return
}

func (s *UserService) Authenticate(ctx context.Context, tokenStr string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	token, err := entity.NewTokenFromEncoded(tokenStr, config.Cfg.HTTP.Secret)
	if err != nil {
		err = apperror.NewError(errors.CodeTokenError, err.Error())
		return
	}
	dbToken, err := s.TokenRepo.GetTokenById(ctx, token.ID)
	if err != nil {
		return
	}
	dbUser, err := s.UserRepo.GetUserById(ctx, dbToken.UserID)
	if err != nil {
		return
	}
	return dbUser, nil
}

func (s *UserService) GetUser(ctx context.Context, e *entity.User) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.UserRepo.GetUserById(ctx, e.ID)
}

func (s *UserService) GetUserList(ctx context.Context, query *repository.QueryParams) (res []*entity.User, total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.UserRepo.GetUserList(ctx, query)
}

func (s *UserService) CreateUser(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	err = s.UserRepo.CheckUserNameAndEmailExist(ctx, e.Username, e.Email)
	if err != nil {
		return
	}

	user := &entity.User{
		ID: utils.NewID(),
	}
	err = user.OnUserCreated(ctx, e, time.Now())
	if err != nil {
		return
	}
	err = s.UserRepo.SaveUser(ctx, user)
	if err != nil {
		return
	}
	_, err = s.OtpSvc.GenerateOtp(ctx, user.Email)
	return
}

func (s *UserService) ActiveUser(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	otp, err := s.OtpSvc.VerifyOtp(ctx, &entity.Otp{
		Email: e.Email,
		Code:  e.Otp,
	})
	if err != nil {
		return
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

	err = s.OtpSvc.DeleteOtp(ctx, otp)
	return
}

func (s *UserService) ResetPassword(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	otp, err := s.OtpSvc.VerifyOtp(ctx, &entity.Otp{
		Email: e.Email,
		Code:  e.Otp,
	})
	if err != nil {
		return
	}
	dbUser, err := s.UserRepo.GetUserByEmail(ctx, e.Email)
	if err != nil {
		return
	}
	dbUser.OnUserUpdated(ctx, e, time.Now())
	err = s.UserRepo.UpdateUser(ctx, dbUser)
	if err != nil {
		return
	}

	err = s.OtpSvc.DeleteOtp(ctx, otp)
	return
}

func (s *UserService) UpdateUser(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	dbUser, err := s.UserRepo.GetUserById(ctx, e.ID)
	if err != nil {
		return
	}
	if e.GetUserName() != "" || e.GetEmail() != "" {
		err = s.UserRepo.CheckDuplicateUserNameAndEmail(ctx, dbUser, e.GetUserName(), e.GetEmail())
		if err != nil {
			return nil, err
		}
	}

	err = dbUser.OnUserUpdated(ctx, e, time.Now())
	if err != nil {
		return
	}
	err = s.UserRepo.UpdateUser(ctx, dbUser)
	return
}

func (s *UserService) DeleteUser(ctx context.Context, e *entity.User) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	dbUser, err := s.UserRepo.GetUserById(ctx, e.ID)
	if err != nil {
		return
	}
	err = dbUser.OnUserDeleted(ctx, e, time.Now())
	if err != nil {
		return
	}
	err = s.UserRepo.DeleteUser(ctx, dbUser)
	return
}
