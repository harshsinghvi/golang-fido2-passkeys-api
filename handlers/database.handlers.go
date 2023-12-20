package handlers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateInDatabase(c *gin.Context, db *gorm.DB, value interface{}, args ...models.Args) bool {
	if res := db.Create(value); res.RowsAffected == 0 || res.Error != nil {
		switch code, _ := utils.PgErrorCodeAndMessage(res.Error); code {
		case pgerrcode.UniqueViolation:
			message := utils.ParseArgs(args, "DuplicateMessage", "Duplicate Fields.").(string)
			BadRequest(c, message)
		default:
			log.Printf("Error While Creating in database: %s", res.Error)
			InternalServerError(c)
		}
		return false
	}
	return true
}

func DeleteInDatabaseById(db *gorm.DB, idField string, id interface{}, value interface{}) bool {
	if res := db.Clauses(clause.Returning{}).Where(fmt.Sprintf("%s = ?", idField), id).Delete(value); res.Error != nil {
		log.Printf("Error While Creating in database: %s", res.Error)
		return false
	}
	return true
}

func CreateChallenge(c *gin.Context, db *gorm.DB, data map[string]interface{}, passkey models.Passkey, args ...models.Args) (bool, models.Challenge) {
	challengeString, challenge, err := utils.CreateChallenge(passkey.PublicKey)
	challenge.UserID = passkey.UserID
	challenge.PasskeyID = passkey.ID
	challenge.Status = models.StatusPending
	challenge.Expiry = time.Now().AddDate(0, 0, 10)

	if err != nil {
		InternalServerError(c)
		return false, models.Challenge{}
	}
	if ok := CreateInDatabase(c, db, &challenge, args...); !ok {
		return false, models.Challenge{}
	}

	data["ChallengeString"] = challengeString
	data["ChallengeID"] = challenge.ID
	data["ChallengeExpiry"] = challenge.Expiry

	return true, challenge
}

func MarkVerified(c *gin.Context, db *gorm.DB, value interface{}, idField string, id string, updateField string, updateValue bool) bool {
	if res := db.Model(value).Where(idField+" = ?", id).Update(updateField, updateValue); res.RowsAffected == 0 || res.Error != nil {
		log.Println("Error While updating verified status Reason: ", res.Error)
		BadRequest(c, "Invalid ID or link expired")
		return false
	}
	return true
}

func TxCommit(c *gin.Context, tx *gorm.DB) bool {
	if res := tx.Commit(); res.Error != nil {
		log.Println("Error Comiting Txn, ", res.Error)
		BadRequest(c, MessageBadRequest)
		return false
	}
	return true
}

func GetById(db *gorm.DB, value interface{}, id interface{}) bool {
	if res := db.First(value, "id = ?", id); res.RowsAffected == 0 || res.Error != nil {
		return false
	}
	return true
}

func LogReqToDb(c *gin.Context, db *gorm.DB, reqId uuid.UUID, reqStart time.Time) {
	accessTokenId, isAuthenticated := c.Get("token_id_uuid")
	billingDisable := c.GetBool("BillingDisable")
	user, _ := c.Get("user")
	hostname, _ := os.Hostname()

	accessLog := &models.AccessLog{
		ID:             reqId,
		RequestID:      reqId,
		Path:           c.Request.URL.Path,
		ServerHostname: hostname,
		ResponseSize:   c.Writer.Size(),
		StatusCode:     c.Writer.Status(),
		ClientIP:       c.ClientIP(),
		Method:         c.Request.Method,
		ResponseTime:   time.Since(reqStart).Milliseconds(),
		Billed:         !isAuthenticated || billingDisable,
		RawQuery:       c.Request.URL.RawQuery,
	}

	if isAuthenticated {
		accessLog.TokenID = accessTokenId.(uuid.UUID)
		accessLog.UserID = user.(models.User).ID
	}

	if ok := CreateInDatabase(c, db, accessLog); !ok {
		log.Println("Error in recording log in db")
	}
}

func UpdateField(db *gorm.DB, value interface{}, idField string, id interface{}, updateField string, updateValue bool) bool {
	res := db.Model(value).Where(idField+" = ?", id).Update(updateField, updateValue)
	if res.RowsAffected == 0 {
		return false
	}
	if res.Error != nil {
		log.Println("Error While updating verified status Reason: ", res.Error)
		return false
	}
	return true
}
