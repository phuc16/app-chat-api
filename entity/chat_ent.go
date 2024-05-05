package entity

import "time"

type Conversation struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	ListUser  []string  `bson:"list_user"`
	Users     []*User   `bson:"users,omitempty"`
	Chat      []*Chat   `bson:"chat"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Chat struct {
	ID               string    `json:"id"`
	FromUserId       string    `json:"from"`
	ToConversationId string    `json:"to"`
	Msg              string    `json:"message"`
	Type             string    `json:"type"`
	Timestamp        time.Time `json:"timestamp"`
}

type Message struct {
	Type              string   `json:"type"`
	ListUserInNewChat []string `json:"list_user,omitempty"`
	Chat              Chat     `json:"chat,omitempty"`
}
