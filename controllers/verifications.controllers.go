package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
)

func Verificaion(c *gin.Context) {
	id := c.Param("id")
	code := c.Query("code")

	if id == "" || code == "" {
		handlers.BadRequest(c, handlers.MessageBadRequestInsufficientData)
		return
	}

	var verification models.Verification

	if ok := handlers.GetById(database.DB, &verification, id); !ok {
		handlers.BadRequest(c, handlers.MessageBadRequest)
		return
	}

	// Checks Based on Request and verification
	if verification.Status == models.StatusSuccess {
		handlers.BadRequest(c, handlers.MessageAlreadyVerified)
		return
	}

	if verification.Status == models.StatusFailed || verification.Status == "" || verification.Type == "" {
		handlers.BadRequest(c, handlers.MessageVerificationAlreadyFailed)
		return
	}

	if time.Until(verification.Expiry).Seconds() <= 0 {
		handlers.BadRequest(c, handlers.MessageExpiredVerificationCode)
		return
	}

	if code != verification.Code {
		handlers.BadRequest(c, handlers.MessageInvalidVerificationCode)
		return
	}

	var ok bool = false
	var message string = handlers.MessageBadRequest

	tx := database.DB.Begin()

	switch verification.Type {
	case models.VerificationTypeNewUser:
		ok = handlers.VerifyNewUser(tx, verification)
		ok, message = utils.CheckBoolAndReturnString(ok, "User Verified", "User Verification Failed or User already verified")
	case models.VerificationTypeNewPasskey:
		ok = handlers.VerifyNewPasskey(tx, verification)
		ok, message = utils.CheckBoolAndReturnString(ok, "Passkey Authorised", "Passkey Authorisation Failed or already authorised")
	case models.VerificationTypeDeleteUser:
		ok = handlers.DeleteUser(tx, verification)
		ok, message = utils.CheckBoolAndReturnString(ok, "User and data deleted", "User Deletion Failed or User already deleted")
	}

	if !ok {
		tx.Rollback()
		verification.Status = models.StatusFailed
		database.DB.Save(&verification)
		handlers.BadRequest(c, message)
		return
	}

	if txOk := handlers.TxCommit(c, tx); !txOk {
		return
	}

	// TODO: Send Conformation Mail
	handlers.StatusOK(c, nil, message)
}

func ReVerifyUser(c *gin.Context) {
	email := c.Param("email")

	if ok := utils.IsEmailValid(email); !ok {
		handlers.BadRequest(c, handlers.MessageInvalidEmailAddress)
		return
	}

	var user models.User
	var passkey models.Passkey
	var challenge models.Challenge

	if res := database.DB.Where("email = ?", email).First(&user); res.RowsAffected == 0 || res.Error != nil {
		handlers.BadRequest(c, "Email address not registered. please register.")
		return
	}

	if user.Verified {
		handlers.BadRequest(c, handlers.MessageAlreadyVerified)
		return
	}

	if res := database.DB.Where("user_id = ?", user.ID).First(&passkey); res.RowsAffected == 0 || res.Error != nil {
		handlers.BadRequest(c, handlers.MessageBadRequest)
		return
	}

	if res := database.DB.Where("user_id = ?", user.ID).First(&challenge); res.RowsAffected == 0 || res.Error != nil {
		handlers.BadRequest(c, handlers.MessageBadRequest)
		return
	}

	verification := utils.CreateVerification(user.ID, models.VerificationTypeNewUser)
	verification.UserID = user.ID
	verification.Email = user.Email

	tx := database.DB.Begin()

	if ok := handlers.CreateInDatabase(c, database.DB, &verification); !ok {
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

	handlers.StatusOK(c, nil, "User Verification Email Sent.")
}

func ReVerifyPasskey(c *gin.Context) {
	var passkey models.Passkey
	var publicKeyStr string
	var user models.User

	if publicKeyStr = c.GetHeader("Public-Key"); publicKeyStr == "" {
		handlers.BadRequest(c, "Public-Key Header not found")
		return
	}

	if res := database.DB.Where("public_key = ?", publicKeyStr).Find(&passkey); res.RowsAffected == 0 {
		handlers.BadRequest(c, "Invalid passkey")
		return
	}

	if passkey.Verified {
		handlers.BadRequest(c, "Passkey already Authorised.")
		return
	}

	if res := database.DB.Where("id = ?", passkey.UserID).First(&user); res.RowsAffected == 0 || res.Error != nil {
		handlers.BadRequest(c, handlers.MessageBadRequest)
		return
	}

	verification := utils.CreateVerification(passkey.ID, models.VerificationTypeNewPasskey)
	verification.UserID = user.ID
	verification.Email = user.Email

	tx := database.DB.Begin()

	if ok := handlers.CreateInDatabase(c, database.DB, &verification); !ok {
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

	handlers.StatusOK(c, nil, "Passkey Authorisation Email Sent.")
}
