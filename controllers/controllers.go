package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models/roles"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/event"
)

func NewUser(c *gin.Context) {
	data := map[string]interface{}{}

	body := handlers.ParseBodyStrict(c, "Email", "Name", "PublicKey") // "PrivateKey"
	if body == nil {
		return
	}

	if ok := utils.IsEmailValid(body["Email"].(string)); !ok {
		handlers.BadRequest(c, handlers.MessageInvalidEmailAddress)
		return
	}

	if ok := utils.IsPublicKeyValid(body["PublicKey"].(string)); !ok {
		handlers.BadRequest(c, handlers.MessageInvalidPublicKey)
		return
	}

	// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
	// if ok := crypto.ValidatePublicAndPrivateKeys(body["PrivateKey"].(string), body["PublicKey"].(string)); !ok {
	// 	handlers.BadRequest(c, "Invalid Public / Private Keys")
	// 	return
	// }

	tx := database.DB.Begin()

	user := models.User{
		Name:     body["Name"].(string),
		Email:    body["Email"].(string),
		Verified: false,
		Roles:    pq.StringArray{roles.User},
	}

	if ok := handlers.CreateInDatabase(c, tx, &user, models.Args{"DuplicateMessage": "Email address already in use, Please use different Email address"}); !ok {
		tx.Rollback()
		return
	}

	passkey := models.Passkey{
		UserID:     user.ID,
		Desciption: "Default Key",
		PublicKey:  body["PublicKey"].(string),
		Verified:   false,
	}

	if ok := handlers.CreateInDatabase(c, tx, &passkey, models.Args{"DuplicateMessage": "Public Key already in use, please Generate new keys."}); !ok {
		tx.Rollback()
		return
	}

	ok, _ := handlers.CreateChallenge(c, tx, data, passkey)
	if !ok {
		tx.Rollback()
		return
	}

	// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
	// passkeyPrivateKey := models.PasskeyPrivateKey{
	// 	UserID:     user.ID,
	// 	PasskeyID:  passkey.ID,
	// 	PrivateKey: body["PrivateKey"].(string),
	// }
	// if ok := handlers.CreateInDatabase(c, tx, &passkeyPrivateKey); !ok {
	// 	tx.Rollback()
	// 	return
	// }

	verification := utils.CreateVerification(user.ID, models.VerificationTypeNewUser)
	verification.UserID = user.ID
	verification.Email = user.Email

	if ok := handlers.CreateInDatabase(c, tx, &verification); !ok {
		tx.Rollback()
		return
	}

	emailOk := handlers.SendVerificationMail(tx, verification)

	if ok := handlers.TxCommit(c, tx); !ok {
		return
	}

	event.PostEvent(database.DB, event.NEW_USER, user.ID.String(), user.Email)

	data["User"] = models.User{
		Name:  user.Name,
		Email: user.Email,
	}

	if !emailOk {
		handlers.StatusOK(c, data, "User Created, Verification Mail not sent, please reverify")
		return
	}
	handlers.StatusOK(c, data, "User Created, please complete Registration by verifing Email, please check your inbox for verification instructions")
}

func VerifyChallenge(c *gin.Context) {
	data := map[string]interface{}{}
	var challenge models.Challenge
	var passkey models.Passkey

	body := handlers.ParseBodyStrict(c, "ChallengeID", "ChallengeSignature")
	if body == nil {
		return
	}

	if res := database.DB.Where("id = ?  AND expiry > now()", body["ChallengeID"].(string)).Find(&challenge); res.RowsAffected == 0 {
		handlers.BadRequest(c, "Invalid/Expired ChallengeID")
		return
	}

	if res := database.DB.Where("id = ?", challenge.PasskeyID).Find(&passkey); res.RowsAffected == 0 {
		handlers.InternalServerError(c)
		return
	}

	if time.Until(challenge.Expiry).Seconds() <= 0 || challenge.Status == models.StatusSuccess || challenge.Status == models.StatusFailed {
		handlers.BadRequest(c, "Challenge Verified Failed, Challenge Expired or Challenge already Verified or Failed")
		return
	}

	message, ok := utils.SolveChallenge(challenge)
	if !ok {
		handlers.InternalServerError(c)
		return
	}

	if ok := utils.VerifySignature(passkey.PublicKey, body["ChallengeSignature"].(string), message); !ok {
		// Update Database
		challenge.Status = models.StatusFailed
		database.DB.Save(&challenge)
		handlers.BadRequest(c, "Challenge Verified Failed")
		return
	}

	accessToken := models.AccessToken{
		UserID:      challenge.UserID,
		PasskeyID:   challenge.PasskeyID,
		ChallengeID: challenge.ID,
		Token:       utils.GenerateToken(challenge.ID.String()),
		Disabled:    !passkey.Verified, // passkey.Verified == false // Token must be disabled when the passkey is not verified
		Expiry:      utils.GenerateTokenExpiryDate(),
		Desciption:  "Generated from Passkey",
	}

	if accessToken.Token == "" {
		handlers.InternalServerError(c)
		return
	}

	if ok := handlers.CreateInDatabase(c, database.DB, &accessToken); !ok {
		return
	}

	challenge.Status = models.StatusSuccess
	database.DB.Save(&challenge)

	data["TokenExpiry"] = accessToken.Expiry
	data["Token"] = accessToken.Token

	handlers.StatusOK(c, data, "Challenge Verification Successful")
}

func RequestChallenge(c *gin.Context) {
	publicKeyStr := c.GetHeader("Public-Key")
	passkeyId := c.Param("passkey")

	var querry *gorm.DB

	if publicKeyStr != "" {
		querry = database.DB.Where("public_key = ?", publicKeyStr)
	} else if passkeyId != "" {
		querry = database.DB.Where("id = ?", passkeyId)
	} else {
		handlers.BadRequest(c, "Public-Key Header or /:passkey-id not found")
		return
	}

	var passkey models.Passkey
	data := map[string]interface{}{}

	if res := querry.Find(&passkey); res.RowsAffected == 0 || res.Error != nil {
		handlers.BadRequest(c, "Invalid passkey")
		return
	}

	if !passkey.Verified {
		handlers.BadRequest(c, "Passkey Not Authorised, please authorise before using.")
		return
	}

	if ok, _ := handlers.CreateChallenge(c, database.DB, data, passkey); !ok {
		handlers.InternalServerError(c)
		return
	}

	handlers.StatusOK(c, data, "Challenge Created, Verify to login")
}

func RegistereNewPasskey(c *gin.Context) {
	body := handlers.ParseBodyStrict(c, "Email", "PublicKey", "Desciption")
	if body == nil {
		return
	}

	tx := database.DB.Begin()

	var user models.User

	if res := tx.Where("email = ?", body["Email"]).Find(&user); res.RowsAffected == 0 || res.Error != nil {
		if !user.Verified {
			handlers.BadRequest(c, "Email Not verified please verify")
			return
		}
		handlers.BadRequest(c, "Email address not found. Please check Email address or register new user.")
		return
	}

	if !user.Verified {
		handlers.BadRequest(c, "User not verified, please check your inbox for instructions")
		return
	}

	passkey := models.Passkey{
		UserID:     user.ID,
		PublicKey:  body["PublicKey"].(string),
		Desciption: body["Desciption"].(string),
		Verified:   false,
	}

	if ok := handlers.CreateInDatabase(c, tx, &passkey, models.Args{"DuplicateMessage": "Public Key already in use, please Generate new keys."}); !ok {
		tx.Rollback()
		return
	}

	verification := utils.CreateVerification(passkey.ID, models.VerificationTypeNewPasskey)
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

	handlers.StatusOK(c, nil, "Passkey Added. Proceed to verification. check your email for verification code or verification link.")
}
