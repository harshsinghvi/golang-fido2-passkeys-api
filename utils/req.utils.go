package utils

import (
	"github.com/mitchellh/mapstructure"
	"log"
)

// INFO Should not use this utility low security
func BindBody(body map[string]interface{}, data interface{}) bool {
	if err := mapstructure.Decode(body, data); err != nil {
		log.Println("Error Decode body, Error: ", err)
		return false
	}
	return true
}
