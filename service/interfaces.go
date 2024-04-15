package service

import (
	"app/entity"
	"app/repository"
	"context"
)

//go:generate mockgen -source $GOFILE -destination ../mocks/$GOPACKAGE/mock_$GOFILE -package mocks

type IUserRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	SaveUser(ctx context.Context, user *entity.User) error
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUserName(ctx context.Context, username string) (*entity.User, error)
	GetUserByUserNameOrEmail(ctx context.Context, username string, email string) (*entity.User, error)
	GetInactiveUser(ctx context.Context, email string) (res *entity.User, err error)
	CheckUserNameAndEmailExist(ctx context.Context, username string, email string) (err error)
	CheckDuplicateUserNameAndEmail(ctx context.Context, user *entity.User, username string, email string) (err error)
	GetUserList(ctx context.Context, params *repository.QueryParams) ([]*entity.User, int64, error)
	GetAllUsers(ctx context.Context) (res []*entity.User, err error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, user *entity.User) error
	CountUser(ctx context.Context) (total int64, err error)
	AddFriendRequest(ctx context.Context, user *entity.User, friend *entity.User) (err error)
	RemoveFriendRequest(ctx context.Context, user *entity.User, friend *entity.User) (err error)
	AddFriend(ctx context.Context, user *entity.User, friend *entity.User) (err error)
	RemoveFriend(ctx context.Context, user *entity.User, friend *entity.User) (err error)
}
type ITokenRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	CreateToken(ctx context.Context, token *entity.Token) error
	GetTokenById(ctx context.Context, id string) (*entity.Token, error)
	GetTokenList(ctx context.Context, params repository.QueryParams) ([]*entity.Token, int64, error)
	UpdateToken(ctx context.Context, token *entity.Token) error
	DeleteToken(ctx context.Context, token *entity.Token) error
}

type IOtpRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	SaveOtp(ctx context.Context, otp *entity.Otp) (err error)
	GetOtp(ctx context.Context, otp *entity.Otp) (res *entity.Otp, err error)
	DeleteOtp(ctx context.Context, otp *entity.Otp) (err error)
}

type IOtpSvc interface {
	GenerateOtp(ctx context.Context, email string) (res entity.Otp, err error)
	VerifyOtp(ctx context.Context, e *entity.Otp) (res *entity.Otp, err error)
	DeleteOtp(ctx context.Context, e *entity.Otp) (err error)
}
