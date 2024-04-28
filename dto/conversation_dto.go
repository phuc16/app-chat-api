package dto

import (
	"app/entity"
	"time"
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
