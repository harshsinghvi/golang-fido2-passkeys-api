package controllers

import (
	"github.com/gin-gonic/gin"
	"harshsinghvi/golang-fido2-passkeys-api/database"
	"harshsinghvi/golang-fido2-passkeys-api/handlers"
	"harshsinghvi/golang-fido2-passkeys-api/models"
	"log"
	"time"
)

func AuthMidlweare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		var accessToken models.AccessToken

		if token == "" {
			handlers.BadRequest(c, "token Not found in Headers.")
			return
		}

		if res := database.DB.Where("token = ?", token).Find(&accessToken); res.RowsAffected == 0 && res.Error != nil {
			if res.Error != nil {
				log.Println("Error in querring auth token, Reason :", res.Error.Error())
			}
			handlers.UnauthorisedRequest(c)
			return
		}

		if time.Until(accessToken.Expiry).Seconds() <= 0 {
			handlers.UnauthorisedRequest(c)
			return
		}

		c.Set("token", accessToken.Token)
		c.Set("user_id", accessToken.UserID)
		c.Next()
	}
}
