package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"log"
	"strings"
)

func ParseBody(c *gin.Context, strict bool, keys []string) map[string]interface{} {
	data := map[string]interface{}{}
	body := map[string]interface{}{}

	if err := c.Bind(&data); err != nil {
		log.Printf("Error Binding request body: %s", err.Error())
		BadRequest(c, MessageInvalidBody)
		return nil
	}

	for _, key := range keys {
		val, ok := data[key]
		if !ok {
			if strict {
				BadRequest(c, fmt.Sprintf("No %s found in request body.", strings.Join([]string(keys), " or ")))
				return nil
			}
		} else {
			body[key] = val
		}
	}

	return body
}

func ParseBodyStrict(c *gin.Context, keys ...string) map[string]interface{} {
	return ParseBody(c, true, keys)
}

func ParseBodyNonStrict(c *gin.Context, keys ...string) map[string]interface{} {
	return ParseBody(c, false, keys)
}

func ParseArgs(args []models.Args, key string, defaultValue interface{}) interface{} {
	if len(args) == 0 {
		return defaultValue
	}

	if val, ok := args[0][key]; ok {
		return val
	}

	return defaultValue
}
