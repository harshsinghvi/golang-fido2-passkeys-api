package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"

	// WIP: TO remove
	// "github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	// "github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
)

func DeleteController(_DataEntity interface{}, args ...models.Args) gin.HandlerFunc {
	defaultMessageValue := fmt.Sprintf("DELETE %s", helpers.GetStructName(_DataEntity))
	_Message := helpers.ParseArgs(args, "Message", defaultMessageValue).(string)
	_SelfResource := helpers.ParseArgs(args, "SelfResource", false).(bool)
	_SelfResourceField := helpers.ParseArgs(args, "SelfResourceField", "user_id").(string)

	return func(c *gin.Context) {
		entityId := c.Param("id")

		querry := database.DB.Where("id = ?", entityId)

		if _SelfResource {
			userId, _ := c.Get("user_id")
			querry = querry.Where(fmt.Sprintf("%s = ?", _SelfResourceField), userId)
		}

		if res := querry.Delete(_DataEntity); res.RowsAffected == 0 || res.Error != nil {
			helpers.BadRequest(c, "Unable to Delete")
			return
		}

		helpers.StatusOK(c, nil, _Message)
	}
}
