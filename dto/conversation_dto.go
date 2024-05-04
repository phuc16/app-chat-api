package dto

import (
	"app/entity"
	"app/errors"
	"app/pkg/apperror"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type ConversationBasicInfoResp struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r ConversationBasicInfoResp) FromConversation(e *entity.Conversation) *ConversationBasicInfoResp {
	return &ConversationBasicInfoResp{
		ID:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

type ConversationInfoResp struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Users     []*UserInfoResp `json:"users"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func fromUserList(userList []*entity.User) (userInfoList []*UserInfoResp) {
	for _, v := range userList {
		userInfoList = append(userInfoList, UserInfoResp{}.FromUser(v))
	}
	return
}

func (r ConversationInfoResp) FromConversation(e *entity.Conversation) *ConversationInfoResp {
	return &ConversationInfoResp{
		ID:        e.ID,
		Name:      e.Name,
		Users:     fromUserList(e.Users),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

type ChatResp struct {
	ID               string    `json:"id"`
	FromUserId       string    `json:"from"`
	ToConversationId string    `json:"to"`
	Msg              string    `json:"message"`
	Type             string    `json:"type"`
	Timestamp        time.Time `json:"timestamp"`
}

func (r ChatResp) FromChat(e *entity.Chat) *ChatResp {
	return &ChatResp{
		ID:               e.ID,
		FromUserId:       e.FromUserId,
		ToConversationId: e.ToConversationId,
		Msg:              e.Msg,
		Type:             e.Type,
		Timestamp:        e.Timestamp,
	}
}

type ChatListResp struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     []*ChatResp `json:"list"`
}

type MessageUpdateReq struct {
	ConversationID string `json:"-"`
	ChatID         string `json:"chat_id"`
	Msg            string `json:"msg"`
}

func (r MessageUpdateReq) Bind(ctx *gin.Context) (*MessageUpdateReq, error) {
	err := ctx.ShouldBindJSON(&r)
	if err != nil {
		return nil, apperror.NewError(errors.CodeUnknownError, validationErrorToText(err))
	}
	return &r, nil
}

func (r MessageUpdateReq) Validate() (err error) {
	return
}

func (r MessageUpdateReq) ToConversation(ctx context.Context) (res *entity.Conversation) {
	res = &entity.Conversation{
		ID: r.ConversationID,
		Chat: []*entity.Chat{
			{
				ID:  r.ChatID,
				Msg: r.Msg,
			},
		},
	}
	return res
}
