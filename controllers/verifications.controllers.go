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
	var verification models.Verification

	if id == "" || code == "" {
		handlers.BadRequest(c, handlers.MessageBadRequestInsufficientData)
		return
	}

	if ok := handlers.GetById(database.DB, &verification, id); !ok {
		handlers.BadRequest(c, handlers.MessageBadRequest)
		return
	}

	// Checks Based on Request and verification
	if verification.Status == models.StatusSuccess {
		handlers.BadRequest(c, handlers.MessageAlreadyVerified)
		return
	}

	if verification.Status == models.StatusFailed {
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

	verification.Status = models.StatusSuccess

	// INFO: Update Database
	var user models.User
	var passkey models.Passkey
	var accessToken models.AccessToken

	tx := database.DB.Begin()

	if verification.UserID != models.NilUUID {
		if ok := handlers.MarkVerified(c, tx, &user, "id", verification.UserID.String(), "verified", true); !ok {
			tx.Rollback()
			return
		}
	}

	if verification.PasskeyID != models.NilUUID {
		if ok := handlers.MarkVerified(c, tx, &passkey, "id", verification.PasskeyID.String(), "verified", true); !ok {
			tx.Rollback()
			return
		}
	}

	if verification.TokenID != models.NilUUID {
		if ok := handlers.MarkVerified(c, tx, &accessToken, "id", verification.TokenID.String(), "disabled", false); !ok {
			tx.Rollback()
			return
		}
	}

	if verification.ChallengeID != models.NilUUID {
		if ok := handlers.MarkVerified(c, tx, &accessToken, "challenge_id", verification.ChallengeID.String(), "disabled", false); !ok {
			tx.Rollback()
			return
		}
	}

	if res := tx.Save(&verification); res.RowsAffected == 0 || res.Error != nil {
		tx.Rollback()
		handlers.BadRequest(c, handlers.MessageBadRequest)
		return
	}

	if ok := handlers.TxCommit(c, tx); !ok {
		return
	}

	if verification.UserID != models.NilUUID {
		handlers.StatusOK(c, nil, handlers.MessageUserVerificationSuccess)
		return
	}

	handlers.StatusOK(c, nil, handlers.MessagePasskeyVerificationSuccess)
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

	verification := models.Verification{
		UserID:      user.ID,
		PasskeyID:   passkey.ID,
		ChallengeID: challenge.ID,
		Status:      models.StatusPending,
		Expiry:      time.Now().AddDate(0, 0, 1),
		Code:        utils.GenerateCode(),
		Email:       user.Email,
	}

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

	verification := models.Verification{
		PasskeyID: passkey.ID,
		Status:    models.StatusPending,
		Expiry:    time.Now().AddDate(0, 0, 1),
		Code:      utils.GenerateCode(),
		Email:     user.Email,
	}

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
