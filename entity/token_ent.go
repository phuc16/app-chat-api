package entity

import (
	"github.com/dgrijalva/jwt-go"
)

const (
	AccessTokenType = "access_token"
)

type Token struct {
	ID       string `json:"id" bson:"id"`
	UserID   string `json:"user_id" bson:"user_id"`
	Name     string `json:"name" bson:"name"`
	UserName string `json:"user_name" bson:"user_name"`
	Email    string `json:"email" bson:"email"`
	Type     string `json:"-" bson:"type"`
	jwt.StandardClaims
}

func NewTokenFromEncoded(tokenStr string, secret string) (*Token, error) {
	token := Token{}
	_, err := jwt.ParseWithClaims(tokenStr, &token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (t *Token) Valid() error {
	return t.StandardClaims.Valid()
}

func (t *Token) SignedToken(secret string) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	signedToken, _ := jwtToken.SignedString([]byte(secret))
	return signedToken
}
