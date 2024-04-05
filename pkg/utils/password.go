package utils

import (
	"app/config"
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hasher.Write([]byte(config.Cfg.HTTP.Secret))
	hashedPassword := hasher.Sum(nil)
	return hex.EncodeToString(hashedPassword)
}

func VerifyPassword(password string, hashedPassword string) bool {
	return HashPassword(password) == hashedPassword
}
