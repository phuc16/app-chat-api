package utils

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

func StringPointer(s string) *string {
	return &s
}

func BooleanPointer(s bool) *bool {
	return &s
}

func Random(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func NewID() string {
	return uuid.NewString()
}

func ContainsString(list []string, str string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}
