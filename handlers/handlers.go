package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models/roles"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/pagination"
	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
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

func StatusOKPag(c *gin.Context, data interface{}, pag pagination.Pagination, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    message,
		"data":       data,
		"pagination": pag.Validate(),
	})
	c.Abort()
}

func UnauthorisedRequest(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  http.StatusUnauthorized,
		"message": "Unauthorised request token invalid/expired/disabled/insufficient roles/permissions.",
	})
	c.Abort()
}

func ParseBody(c *gin.Context, strict bool, keys []string) map[string]interface{} {
	data := map[string]interface{}{}
	body := map[string]interface{}{}

	if err := c.Bind(&data); err != nil {
		log.Printf("Error Binding request body: %s", err.Error())
		BadRequest(c, MessageInvalidBody)
		return nil
	}

	for _, key := range keys {
		val, ok := data[key]
		if !ok {
			if strict {
				BadRequest(c, fmt.Sprintf("No %s found in request body.", strings.Join([]string(keys), " or ")))
				return nil
			}
		} else {
			body[key] = val
		}
	}

	return body
}

func ParseBodyStrict(c *gin.Context, keys ...string) map[string]interface{} {
	return ParseBody(c, true, keys)
}

func ParseBodyNonStrict(c *gin.Context, keys ...string) map[string]interface{} {
	return ParseBody(c, false, keys)
}

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

func HealthHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "OK"})
}

func ExternalRedirect(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, url)
	}
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

func GetById(db *gorm.DB, value interface{}, id string) bool {
	if res := db.First(value, "id = ?", id); res.RowsAffected == 0 || res.Error != nil {
		return false
	}
	return true
}

func LogReqToDb(c *gin.Context, db *gorm.DB, reqId uuid.UUID, reqStart time.Time) {
	accessTokenId, isAuthenticated := c.Get("token_id_uuid")
	billingDisable := c.GetBool("BillingDisable")
	userId, _ := c.Get("user_id_uuid")
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
	}

	if isAuthenticated {
		accessLog.TokenID = accessTokenId.(uuid.UUID)
		accessLog.UserID = userId.(uuid.UUID)
	}

	if ok := CreateInDatabase(c, db, accessLog); !ok {
		log.Println("Error in recording log in db")
	}
}

// TODO Test and usage pending
func CheckForSelfResource(c *gin.Context, value interface{}) bool {
	userId, oKa := c.Get("user_id")
	userRoles, oKb := c.Get("user_roles")

	if !oKa || !oKb {
		UnauthorisedRequest(c)
		return false
	}

	if ok := roles.CheckRoles([]string{roles.Admin, roles.SuperAdmin}, userRoles.([]string)); ok {
		return true
	}

	switch entity := value.(type) {
	case models.User:
		return userId.(string) == entity.ID.String()
	case models.Passkey:
		return userId.(string) == entity.UserID.String()
	case models.Challenge:
		return userId.(string) == entity.UserID.String()
	case models.AccessToken:
		return userId.(string) == entity.UserID.String()
	case models.Verification:
		return userId.(string) == entity.UserID.String()
	default:
		UnauthorisedRequest(c)
		return false
	}
}
