package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"harshsinghvi/golang-fido2-passkeys-api/database"
	"harshsinghvi/golang-fido2-passkeys-api/handlers"
	"harshsinghvi/golang-fido2-passkeys-api/models"
	"harshsinghvi/golang-fido2-passkeys-api/utils"
	"log"
	"time"
)

func GetControler(model interface{}) gin.HandlerFunc {
	// entityName := reflect.TypeOf(&model{})
	// search field =
	log.Println()
	return func(c *gin.Context) {
		var users models.Users
		database.DB.Find(&users)
	}
}

func NewUser(c *gin.Context) {
	var user models.User
	var passkey models.Passkey
	// var passkeyPrivateKey models.PasskeyPrivateKey

	data := map[string]interface{}{}

	// body := handlers.ParseBody(c, []string{"Email", "Name", "PrivateKey", "PublicKey"})
	body := handlers.ParseBody(c, []string{"Email", "Name", "PublicKey"})
	if body == nil {
		handlers.InternalServerError(c)
		return
	}

	if err := mapstructure.Decode(body, &user); err != nil {
		handlers.InternalServerError(c)
		return
	}

	// if ok := crypto.ValidatePublicAndPrivateKeys(body["PrivateKey"].(string), body["PublicKey"].(string)); !ok {
	// 	handlers.BadRequest(c, "Invalid Public / Private Keys")
	// 	return
	// }

	if ok := handlers.CreateInDatabase(c, &user); !ok {
		return
	}

	passkey.UserID = user.ID
	passkey.Desciption = "Default Key"
	passkey.PublicKey, _ = body["PublicKey"].(string)

	if ok := handlers.CreateInDatabase(c, &passkey); !ok {
		return
	}

	// passkeyPrivateKey.UserID = user.ID
	// passkeyPrivateKey.PasskeyID = passkey.ID
	// passkeyPrivateKey.PrivateKey, _ = body["PrivateKey"].(string)

	// if ok := handlers.CreateInDatabase(c, &passkeyPrivateKey); !ok {
	// 	return
	// }

	if ok := handlers.CreateChallenge(c, data, passkey); !ok {
		handlers.InternalServerError(c)
		return
	}

	data["PasskeyID"] = passkey.ID
	data["User"] = user
	handlers.StatusOK(c, data, "User Created, please complete Registration by completing challenge")
}

func VerifyChallenge(c *gin.Context) {
	data := map[string]interface{}{}

	body := handlers.ParseBody(c, []string{"ChallengeID", "ChallengeSignature"})
	if body == nil {
		return
	}

	var challenge models.Challenge
	var passkey models.Passkey

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

	if ok := handlers.VerifySignature(passkey.PublicKey, body["ChallengeSignature"].(string), message); !ok {
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

	if ok := handlers.CreateInDatabase(c, &accessToken); !ok {
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

	if ok := handlers.CreateChallenge(c, data, passkey); !ok {
		handlers.InternalServerError(c)
		return
	}

	handlers.StatusOK(c, data, "Challenge Created, Verify to login")
}

func RequestChallengeUsingPublicKey(c *gin.Context) {
	publicKey := c.GetHeader("Public-Key")

	if publicKey == "" {
		handlers.BadRequest(c, "Public-Key Header not found")
		return
	}

	data := map[string]interface{}{}
	var passkey models.Passkey

	if res := database.DB.Where("public_key = ?", publicKey).Find(&passkey); res.RowsAffected == 0 {
		handlers.BadRequest(c, "Invalid passkey")
		return
	}

	if ok := handlers.CreateChallenge(c, data, passkey); !ok {
		handlers.InternalServerError(c)
		return
	}

	handlers.StatusOK(c, data, "Challenge Created, Verify to login")
}
