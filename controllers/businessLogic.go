package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
)

func GetMe(c *gin.Context) {
	data := map[string]interface{}{}
	data["Token"], _ = c.Get("token")
	data["UserID"], _ = c.Get("user_id")
	handlers.StatusOK(c, data, "Request Authenticated")
}
