package helpers

import (
	"time"
)

func TimeNow(args ...interface{}) interface{} {
	return time.Now()
}
