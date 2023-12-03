package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/jackc/pgerrcode"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	"strings"
	"time"
)

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": message,
	})
	c.Abort()
}

func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Internal server error, Something went wrong !",
	})
	c.Abort()
}

func StatusOK(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": message,
		"data":    data,
	})
	c.Abort()
}

func UnauthorisedRequest(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  http.StatusUnauthorized,
		"message": "Unauthorised request token invalid/expired/disabled.",
	})
	c.Abort()
}

func ParseBody(c *gin.Context, keys []string) map[string]interface{} {
	body := map[string]interface{}{}
	if err := c.Bind(&body); err != nil {
		log.Printf("Error Binding request body: %s", err.Error())
		BadRequest(c, "Invalid body")
		return nil
	}

	for _, key := range keys {
		if _, ok := body[key]; !ok {
			BadRequest(c, fmt.Sprintf("No %s found in request body.", strings.Join([]string(keys), " or ")))
			return nil
		}
	}

	return body
}

func ParseBodyAndBind(c *gin.Context, keys []string, data interface{}) bool {
	body := ParseBody(c, keys)

	if body == nil {
		InternalServerError(c)
		return false
	}

	if err := mapstructure.Decode(body, &data); err != nil {
		InternalServerError(c)
		return false
	}
	return true
}

func CreateInDatabase(c *gin.Context, value interface{}) bool {
	if res := database.DB.Create(value); res.RowsAffected == 0 || res.Error != nil {
		switch code, _ := utils.PgErrorCodeAndMessage(res.Error); code {
		case pgerrcode.UniqueViolation:
			BadRequest(c, "Duplicate Fields")
			return false
		default:
			log.Printf("Error While Creating in database: %s", res.Error.Error())
			InternalServerError(c)
			return false
		}
	}
	return true
}

func VerifySignature(publicKeyStr string, signatureStr string, message string) bool {
	publicKey, err := crypto.ParsePublicKey(publicKeyStr)
	if err != nil {
		log.Println("Error While parsing public key from db", err)
	}

	signature, err := base64.StdEncoding.DecodeString(signatureStr)
	if err != nil {
		log.Println("Error Parsing signature :", err)
		return false
	}

	err = crypto.VerifySignature(signature, publicKey, message)
	if err != nil {
		fmt.Println("Signature verification failed:", err)
		return false
	}

	return true
}

func CreateChallenge(c *gin.Context, data map[string]interface{}, passkey models.Passkey) bool {
	challengeString, challenge, err := utils.CreateChallenge(passkey.PublicKey)

	challenge.UserID = passkey.UserID
	challenge.PasskeyID = passkey.ID
	challenge.Status = "PENDING"
	challenge.Expiry = time.Now().AddDate(0, 0, 10)

	if ok := CreateInDatabase(c, &challenge); !ok || err != nil {
		return false
	}

	data["ChallengeString"] = challengeString
	data["ChallengeID"] = challenge.ID
	data["ChallengeExpiry"] = challenge.Expiry

	return true
}

func HealthHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func ExternalRedirect(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, url)
	}
}
