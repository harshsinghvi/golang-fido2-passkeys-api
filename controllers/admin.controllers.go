package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
)

func VerifyPasskey(c *gin.Context) {
	id := c.Param("id")
	var passkey models.Passkey

	if ok := handlers.MarkVerified(c, database.DB, &passkey, "id", id, "verified", true); !ok {
		return
	}

	handlers.StatusOK(c, nil, "Passkey Verified")
}

func VerifyUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	var passkey models.Passkey
	var accessToken models.AccessToken

	tx := database.DB.Begin()

	tx.First(&user, "id = ?", id)

	if user.Verified {
		handlers.BadRequest(c, "User Already Verified.")
		return
	}

	if ok := handlers.MarkVerified(c, tx, &user, "id", id, "verified", true); !ok {
		tx.Rollback()
		return
	}

	if ok := handlers.MarkVerified(c, tx, &passkey, "user_id", id, "verified", true); !ok {
		tx.Rollback()
		return
	}

	if ok := handlers.MarkVerified(c, tx, &accessToken, "user_id", id, "disabled", false); !ok {
		tx.Rollback()
		return
	}

	if ok := handlers.TxCommit(c, tx); !ok {
		return
	}

	handlers.StatusOK(c, nil, "User Verified.")
}
