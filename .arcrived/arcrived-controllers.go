package arcrived

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models/roles"
)

func RequestChallengeByID(c *gin.Context) {
	passkeyId := c.Param("passkey")
	data := map[string]interface{}{}
	var passkey models.Passkey

	if res := database.DB.Where("id = ?", passkeyId).Find(&passkey); res.RowsAffected == 0 {
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

func RequestChallengeUsingPublicKey(c *gin.Context) {
	data := map[string]interface{}{}
	var passkey models.Passkey
	publicKeyStr := c.GetHeader("Public-Key")

	if publicKeyStr == "" {
		handlers.BadRequest(c, "Public-Key Header not found")
		return
	}

	if res := database.DB.Where("public_key = ?", publicKeyStr).Find(&passkey); res.RowsAffected == 0 {
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

// INFO Test and usage pending
func CheckForSelfResource(c *gin.Context, value interface{}) bool {
	userId, oKa := c.Get("user_id")
	userRoles, oKb := c.Get("user_roles")

	if !oKa || !oKb {
		handlers.UnauthorisedRequest(c)
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
		handlers.UnauthorisedRequest(c)
		return false
	}
}
