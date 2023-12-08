package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
)

func GetMe(c *gin.Context) {
	data := map[string]interface{}{}
	// data["Token"], _ = c.Get("token")
	userId, _ := c.Get("user_id")
	var user models.User
	if ok := handlers.GetById(c, database.DB, &user, userId.(string)); !ok {
		handlers.BadRequest(c, "No Data Found")
		return
	}
	data["User"] = user
	handlers.StatusOK(c, data, "Request Authenticated")
}
