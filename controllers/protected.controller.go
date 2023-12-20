package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/event"
)

func GetMe(c *gin.Context) {
	data := map[string]interface{}{}
	data["Token"], _ = c.Get("token")
	data["User"], _ = c.Get("user")
	handlers.StatusOK(c, data, "Request Authenticated")
}

// EMAIL
func DeleteUserAndData(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(models.User)

	tx := database.DB.Begin()

	verification := utils.CreateVerification(user.ID, models.VerificationTypeDeleteUser)
	verification.UserID = user.ID
	verification.Email = user.Email

	if ok := handlers.CreateInDatabase(c, tx, &verification); !ok {
		tx.Rollback()
		return
	}

	if ok := handlers.SendVerificationMail(tx, verification); !ok {
		tx.Rollback()
		handlers.BadRequest(c, handlers.MessageErrorWhileSendingEmail)
		return
	}

	if ok := handlers.TxCommit(c, tx); !ok {
		return
	}

	event.PostEvent(database.DB, event.DELETE_USER, user.ID.String(), user.Email)
	handlers.StatusOK(c, nil, "User Deletion requested please check your email inbox.")
}
