package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/event"
)

func GetMe(c *gin.Context) {
	data := map[string]interface{}{}
	// data["Token"], _ = c.Get("token")
	userId, _ := c.Get("user_id")
	var user models.User
	if ok := handlers.GetById(database.DB, &user, userId); !ok {
		handlers.BadRequest(c, "No Data Found")
		return
	}
	data["User"] = user
	handlers.StatusOK(c, data, "Request Authenticated")
}

// EMAIL
func DeleteUserAndData(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(models.User)

	tx := database.DB.Begin()

	verification := utils.CreateVerification(user.ID, models.VerificationTypeDeleteUser)
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

	event.PostEvent(database.DB, event.DELETE_USER, user.ID.String(), user.Email)
	handlers.StatusOK(c, nil, "User Deletion requested please check your email inbox.")
}

// func VerifyDeleteUser(c *gin.Context) {
// 	userId, _ := c.Get("user_id")

// 	tx := database.DB.Begin()

// 	if ok := handlers.DeleteInDatabaseById(tx, "id", userId, &[]models.User{}); !ok {
// 		tx.Rollback()
// 		return
// 	}

// 	if ok := handlers.DeleteInDatabaseById(c, tx, "user_id", userId, &[]models.Passkey{}); !ok {
// 		tx.Rollback()
// 		return
// 	}

// 	// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
// 	// if ok := handlers.DeleteInDatabaseById(c, tx, "user_id", userId, &[]models.PasskeyPrivateKey{}); !ok {
// 	// 	tx.Rollback()
// 	// 	return
// 	// }

// 	if ok := handlers.DeleteInDatabaseById(c, tx, "user_id", userId, &[]models.Challenge{}); !ok {
// 		tx.Rollback()
// 		return
// 	}

// 	if ok := handlers.DeleteInDatabaseById(c, tx, "user_id", userId, &[]models.AccessToken{}); !ok {
// 		tx.Rollback()
// 		return
// 	}

// 	if ok := handlers.DeleteInDatabaseById(c, tx, "user_id", userId, &[]models.AccessLog{}); !ok {
// 		tx.Rollback()
// 		return
// 	}

// 	if ok := handlers.DeleteInDatabaseById(c, tx, "user_id", userId, &[]models.Verification{}); !ok {
// 		tx.Rollback()
// 		return
// 	}

// 	if ok := handlers.TxCommit(c, tx); !ok {
// 		return
// 	}

// 	handlers.StatusOK(c, nil, "Delete OK")
// }
