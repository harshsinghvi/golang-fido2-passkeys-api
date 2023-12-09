package utils

import (
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"time"
)

func TimeNow(args ...interface{}) interface{} {
	return time.Now()
}

func TimeNowAfterDays(days int) models.GenFunc {
	return func(args ...interface{}) interface{} {
		return time.Now().AddDate(0, 0, days)
	}
}
