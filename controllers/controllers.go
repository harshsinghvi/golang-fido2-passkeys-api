package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"time"
)

// func GetControler(model interface{}) gin.HandlerFunc {
// 	// entityName := reflect.TypeOf(&model{})
// 	// search field =
// 	log.Println()
// 	return func(c *gin.Context) {
// 		var users models.Users
// 		database.DB.Find(&users)
// 	}
// }

func NewUser(c *gin.Context) {
	var user models.User
	var passkey models.Passkey
	data := map[string]interface{}{}

	// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
	// var passkeyPrivateKey models.PasskeyPrivateKey
	// body := handlers.ParseBody(c, []string{"Email", "Name", "PrivateKey", "PublicKey"})
	body := handlers.ParseBody(c, []string{"Email", "Name", "PublicKey"})
	if body == nil {
		return
	}

	if ok := utils.BindBody(body, &user); !ok {
		handlers.BadRequest(c, "Invalid body")
		return
	}

	// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
	// if ok := crypto.ValidatePublicAndPrivateKeys(body["PrivateKey"].(string), body["PublicKey"].(string)); !ok {
	// 	handlers.BadRequest(c, "Invalid Public / Private Keys")
	// 	return
	// }

	tx := database.DB.Begin()

	if ok := handlers.CreateInDatabase(c, tx, &user, models.Args{"DuplicateMessage": "Email address already in use, Please use different Email address"}); !ok {
		tx.Rollback()
		return
	}

	passkey.UserID = user.ID
	passkey.Desciption = "Default Key"
	passkey.PublicKey, _ = body["PublicKey"].(string)

	if ok := handlers.CreateInDatabase(c, tx, &passkey, models.Args{"DuplicateMessage": "Public Key already in use, please Generate new keys."}); !ok {
		tx.Rollback()
		return
	}

	if ok := handlers.CreateChallenge(c, tx, data, passkey); !ok {
		tx.Rollback()
		return
	}

	// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
	// passkeyPrivateKey.UserID = user.ID
	// passkeyPrivateKey.PasskeyID = passkey.ID
	// passkeyPrivateKey.PrivateKey, _ = body["PrivateKey"].(string)
	// if ok := handlers.CreateInDatabase(c, tx, &passkeyPrivateKey, models.Args{"tx": tx}); !ok {
	// 	tx.Rollback()
	// 	return
	// }

	tx.Commit()

	data["PasskeyID"] = passkey.ID
	data["User"] = user
	handlers.StatusOK(c, data, "User Created, please complete Registration by completing challenge")
}

func VerifyChallenge(c *gin.Context) {
	data := map[string]interface{}{}
	var challenge models.Challenge
	var passkey models.Passkey

	body := handlers.ParseBody(c, []string{"ChallengeID", "ChallengeSignature"})
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

	if time.Until(challenge.Expiry).Seconds() <= 0 || challenge.Status == "SUCCESS" {
		handlers.BadRequest(c, "Challenge Verified Failed, Challenge Expired or Challenge already verified")
		return
	}

	message, ok := utils.SolveChallenge(challenge)
	if !ok {
		handlers.InternalServerError(c)
		return
	}

	if ok := utils.VerifySignature(passkey.PublicKey, body["ChallengeSignature"].(string), message); !ok {
		// Update Database
		challenge.Status = "FAILED"
		database.DB.Save(&challenge)
		handlers.BadRequest(c, "Challenge Verified Failed")
		return
	}

	var accessToken models.AccessToken
	accessToken.UserID = challenge.UserID
	accessToken.PasskeyID = challenge.PasskeyID
	accessToken.ChallengeID = challenge.ID
	accessToken.Token = utils.GenerateToken(challenge.ID.String())
	accessToken.Disabled = false
	accessToken.Expiry = time.Now().AddDate(0, 0, 10)

	if accessToken.Token == "" {
		handlers.InternalServerError(c)
		return
	}

	if ok := handlers.CreateInDatabase(c, database.DB, &accessToken); !ok {
		return
	}

	challenge.Status = "SUCCESS"
	database.DB.Save(&challenge)

	data["TokenExpiry"] = accessToken.Expiry
	data["Token"] = accessToken.Token

	handlers.StatusOK(c, data, "Challenge Verification Successful")
}

func RequestChallenge(c *gin.Context) {
	passkeyId := c.Param("passkey")
	data := map[string]interface{}{}
	var passkey models.Passkey

	if res := database.DB.Where("id = ?", passkeyId).Find(&passkey); res.RowsAffected == 0 {
		handlers.BadRequest(c, "Invalid passkey")
		return
	}

	if ok := handlers.CreateChallenge(c, database.DB, data, passkey); !ok {
		handlers.InternalServerError(c)
		return
	}

	handlers.StatusOK(c, data, "Challenge Created, Verify to login")
}

func RequestChallengeUsingPublicKey(c *gin.Context) {
	data := map[string]interface{}{}
	var passkey models.Passkey
	var publicKeyStr string

	if publicKeyStr = c.GetHeader("Public-Key"); publicKeyStr == "" {
		handlers.BadRequest(c, "Public-Key Header not found")
		return
	}

	if res := database.DB.Where("public_key = ?", publicKeyStr).Find(&passkey); res.RowsAffected == 0 {
		handlers.BadRequest(c, "Invalid passkey")
		return
	}

	if ok := handlers.CreateChallenge(c, database.DB, data, passkey); !ok {
		handlers.InternalServerError(c)
		return
	}

	handlers.StatusOK(c, data, "Challenge Created, Verify to login")
}
