package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"log"
	"time"
)

func AuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken models.AccessToken
		token := c.GetHeader("token")

		if token == "" {
			handlers.BadRequest(c, "token Not found in Headers.")
			return
		}

		res := database.DB.Where("token = ? AND disabled = false AND expiry > now()", token).Find(&accessToken)

		if res.RowsAffected == 0 || res.Error != nil {
			if res.Error != nil {
				log.Println("Error in querring auth token, Reason :", res.Error)
			}
			handlers.UnauthorisedRequest(c)
			return
		}

		c.Set("token", accessToken.Token)
		c.Set("token_id_uuid", accessToken.ID)
		c.Set("user_id", accessToken.UserID.String())
		c.Set("user_id_uuid", accessToken.UserID)
		c.Next()
	}
}

func LoggerMW(args ...models.Args) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := uuid.New()
		reqStart := time.Now()
		c.Set("requestId", reqId)
		c.Writer.Header().Set("X-Request-Id", reqId.String())
		c.Next()
		handlers.LogReqToDb(c, database.DB, reqId, reqStart)
	}
}

func ConfigMW(args ...models.Args) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.ParseAndSet(c, args, "BillingDisable", false)
		c.Next()
	}
}
