package helpers

import (
	"time"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
)

func TimeNow(args ...interface{}) interface{} {
	return time.Now()
}

func TimeNowAfterDays(days int) models.GenerateFunction {
	return func(args ...interface{}) interface{} {
		return time.Now().AddDate(0, 0, days)
	}
}
