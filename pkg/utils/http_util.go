package utils

import (
	"app/pkg/logger"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetBearerAuth(c *gin.Context) (string, bool) {
	auth := c.Request.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = c.Request.FormValue("access_token")
	}
	return token, token != ""
}

func ValidationErrorToText(err error) string {
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, e := range err {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", e.Field())
			case "max":
				return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
			case "min":
				return fmt.Sprintf("%s must be longer than %s", e.Field(), e.Param())
			case "email":
				return fmt.Sprintf("Invalid email format")
			case "len":
				return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
			}
			return fmt.Sprintf("%s is not valid", e.Field())
		}
	}
	return err.Error()
}

func HTTPLogger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	c.Next()
	if raw != "" {
		path = path + "?" + raw
	}
	fields := []zapcore.Field{}
	fields = append(fields, zap.Int("STATUS", c.Writer.Status()))
	fields = append(fields, zap.String("METHOD", c.Request.Method))
	fields = append(fields, zap.String("PATH", path))
	fields = append(fields, zap.Int64("DUR", time.Now().Sub(start).Milliseconds()))
	fields = append(fields, zap.String("IP", c.ClientIP()))
	if value := c.Value("username"); value != nil {
		fields = append(fields, zap.String("username", value.(string)))
	}
	logger.For(c.Request.Context()).Append(fields...).Info("HTTP Request")
}
