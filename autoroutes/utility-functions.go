package autoroutes

import (
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
)

// Utility functions for config generation
func GenerateConstantValue(val interface{}) models.GenerateFunction {
	return func(args ...interface{}) interface{} {
		return val
	}
}