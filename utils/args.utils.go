package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
)

func ParseArgs(args []models.Args, key string, defaultValue interface{}) interface{} {
	if args == nil {
		return defaultValue
	}

	if val, ok := args[0][key]; ok {
		return val
	}
	return defaultValue
}

func ParseAndSet(c *gin.Context, args []models.Args, key string, defaultValue interface{}) {
	value := ParseArgs(args, key, false).(bool)
	c.Set(key, value)
}
