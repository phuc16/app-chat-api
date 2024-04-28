package entity

import "time"

type Conversation struct {
	ID       string   `bson:"id"`
	Name     string   `bson:"name"`
	ListUser []string `bson:"list_user"`
	Chat     []Chat   `bson:"chat"`
}

type Chat struct {
	ID               string    `json:"id"`
	FromUserId       string    `json:"from"`
	ToConversationId string    `json:"to"`
	Msg              string    `json:"message"`
	Timestamp        time.Time `json:"timestamp"`
}

type Message struct {
	Type              string   `json:"type"`
	ListUserInNewChat []string `json:"user,omitempty"`
	Chat              Chat     `json:"chat,omitempty"`
}
