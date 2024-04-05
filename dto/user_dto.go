package dto

import (
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type UserLoginReq struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

func (r UserLoginReq) Bind(ctx *gin.Context) (res *UserLoginReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}
func (r UserLoginReq) Validate() (err error) {
	return
}

func (r UserLoginReq) ToUser(ctx context.Context) (res *entity.User) {
	res = &entity.User{
		Username: r.UsernameOrEmail,
		Email:    r.UsernameOrEmail,
		Password: r.Password,
	}
	return res
}

type UserCreateReq struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"email"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	AvatarUrl string `json:"avatar_url" binding:"required"`
}

func (r UserCreateReq) Bind(ctx *gin.Context) (res *UserCreateReq, err error) {
	err = ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}
func (r UserCreateReq) Validate() (err error) {
	return
}

func (r UserCreateReq) ToUser(ctx context.Context) (res *entity.User) {
	res = &entity.User{
		Username:  r.Username,
		Email:     r.Email,
		Password:  r.Password,
		Name:      r.Name,
		Phone:     r.Phone,
		AvatarUrl: r.AvatarUrl,
	}
	return res
}

type UserUpdateReq struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	AvatarUrl string `json:"avatar_url"`
}

func (r UserUpdateReq) Bind(ctx *gin.Context) (*UserUpdateReq, error) {
	err := ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r UserUpdateReq) Validate() (err error) {
	return
}

func (r UserUpdateReq) ToUser(ctx context.Context) (res *entity.User) {
	res = &entity.User{
		ID:        entity.GetUserFromContext(ctx).ID,
		Username:  r.Username,
		Email:     r.Email,
		Password:  r.Password,
		Name:      r.Name,
		Phone:     r.Phone,
		AvatarUrl: r.AvatarUrl,
	}
	return res
}

type UserDeleteReq struct {
	ID string `json:"id" binding:"required"`
}

func (r UserDeleteReq) Bind(ctx *gin.Context) (*UserDeleteReq, error) {
	err := ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r UserDeleteReq) Validate() (err error) {
	return
}

func (r UserDeleteReq) ToUser(ctx context.Context) (res *entity.User) {
	res = &entity.User{
		ID: r.ID,
	}
	return res
}

type UserResp struct {
	ID              string          `json:"id"`
	Username        string          `json:"username"`
	Email           string          `json:"email"`
	Name            string          `json:"name"`
	Phone           string          `json:"phone"`
	AvatarUrl       string          `json:"avatar_url"`
	Status          string          `json:"status"`
	Friends         []*UserInfoResp `json:"friends"`
	ConversationIds []string        `json:"conversation_ids"`
	LastLoggedIn    time.Time       `json:"last_logged_in"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (r UserResp) FromUser(e *entity.User) *UserResp {
	return &UserResp{
		ID:              e.ID,
		Username:        e.Username,
		Email:           e.Email,
		Name:            e.Name,
		Phone:           e.Phone,
		AvatarUrl:       e.AvatarUrl,
		Status:          e.Status,
		Friends:         fromFriendList(e.Friends),
		ConversationIds: e.ConversationIds,
		LastLoggedIn:    e.LastLoggedIn,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func fromFriendList(friendList []*entity.User) (userInfoList []*UserInfoResp) {
	for _, v := range friendList {
		userInfoList = append(userInfoList, UserInfoResp{}.FromUser(v))
	}
	return
}

type UserListResp struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     []*UserResp `json:"list"`
}

type UserInfoResp struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	AvatarUrl    string    `json:"avatar_url"`
	Status       string    `json:"status"`
	LastLoggedIn time.Time `json:"last_logged_in"`
}

func (r UserInfoResp) FromUser(e *entity.User) *UserInfoResp {
	if e == nil {
		return nil
	}
	return &UserInfoResp{
		ID:           e.ID,
		Username:     e.Username,
		Email:        e.Email,
		Name:         e.Name,
		Phone:        e.Phone,
		AvatarUrl:    e.AvatarUrl,
		Status:       e.Status,
		LastLoggedIn: e.LastLoggedIn,
	}
}
