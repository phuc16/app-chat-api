package entity

import (
	"app/pkg/utils"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

const userKey = "user"
const userNameKey = "username"

const (
	UserStatusOnline  = "online"
	UserStatusOffline = "offline"
	UseStatusAway     = "away"
)

type User struct {
	ID               string     `bson:"id"`
	Username         string     `bson:"username"`
	Email            string     `bson:"email"`
	Password         string     `bson:"password"`
	Name             string     `bson:"name"`
	Phone            string     `bson:"phone"`
	AvatarUrl        string     `bson:"avatar_url"`
	Status           string     `bson:"status"`
	FriendRequestIds []string   `bson:"friend_request_ids"`
	FriendRequests   []*User    `bson:"friend_requests,omitempty"`
	FriendIds        []string   `bson:"friend_ids"`
	Friends          []*User    `bson:"friends,omitempty"`
	ConversationIds  []string   `bson:"conversation_ids"`
	IsActive         bool       `bson:"is_active"`
	LastLoggedIn     time.Time  `bson:"last_logged_in"`
	Otp              string     `bson:"-"`
	CreatedAt        time.Time  `bson:"created_at"`
	UpdatedAt        time.Time  `bson:"updated_at"`
	DeletedAt        *time.Time `bson:"deleted_at,omitempty"`
}

func (e User) GetUserName() string {
	if e.Username != "" {
		return e.Username
	}
	return ""
}

func (e User) GetEmail() string {
	if e.Email != "" {
		return e.Email
	}
	return ""
}

func (e *User) LoggedIn() bool {
	e.LastLoggedIn = time.Now()
	return false
}

func (e *User) OnUserCreated(ctx context.Context, user *User, eventTime time.Time) error {
	e.Username = user.Username
	e.Email = user.Email
	e.Password = utils.HashPassword(user.Password)
	e.Name = user.Name
	e.Phone = user.Phone
	e.AvatarUrl = user.AvatarUrl
	e.Status = UserStatusOffline
	e.FriendRequestIds = make([]string, 0)
	e.FriendIds = make([]string, 0)
	e.ConversationIds = make([]string, 0)
	e.CreatedAt = eventTime
	e.UpdatedAt = eventTime
	return nil
}

func (e *User) OnUserUpdated(ctx context.Context, user *User, eventTime time.Time) error {
	if user.Username != "" {
		e.Username = user.Username
	}
	if user.Email != "" {
		e.Email = user.Email
	}
	if user.Password != "" {
		e.Password = utils.HashPassword(user.Password)
	}
	if user.Name != "" {
		e.Name = user.Name
	}
	if user.Phone != "" {
		e.Phone = user.Phone
	}
	if user.AvatarUrl != "" {
		e.AvatarUrl = user.AvatarUrl
	}
	e.UpdatedAt = eventTime
	return nil
}

func (e *User) OnUserDeleted(ctx context.Context, user *User, eventTime time.Time) error {
	e.DeletedAt = &eventTime
	return nil
}

func (e *User) OnUserActive(ctx context.Context) error {
	e.IsActive = true
	return nil
}

func (e *User) SetToContext(ctx *gin.Context) {
	ctx.Set(userNameKey, e.Username)
	newCtx := context.WithValue(ctx.Request.Context(), userKey, *e)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

func GetUserFromContext(ctx context.Context) User {
	user := User{}
	value := ctx.Value(userKey)
	if value != nil {
		user = value.(User)
	}
	return user
}
