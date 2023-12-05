package utils

import (
	"github.com/mitchellh/mapstructure"
)

func BindBody(body map[string]interface{}, data interface{}) bool {
	if err := mapstructure.Decode(body, &data); err != nil {
		return false
	}
	return true
}
