package entity

import "time"

type Otp struct {
	ID        string    `bson:"id"`
	Email     string    `bson:"email"`
	Code      string    `bson:"code"`
	CreatedAt time.Time `bson:"created_at"`
}
