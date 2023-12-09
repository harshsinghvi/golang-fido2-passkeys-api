package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
)

func DeleteController(_DataEntity interface{}, args ...models.Args) gin.HandlerFunc {
	defaultMessageValue := fmt.Sprintf("DELETE %s", utils.GetStructName(_DataEntity))
	_Message := utils.ParseArgs(args, "Message", defaultMessageValue).(string)
	_SelfResource := utils.ParseArgs(args, "SelfResource", false).(bool)
	_SelfResourceField := utils.ParseArgs(args, "SelfResourceField", "user_id").(string)

	return func(c *gin.Context) {
		entityId := c.Param("id")

		querry := database.DB.Where("id = ?", entityId)

		if _SelfResource {
			userId, _ := c.Get("user_id")
			querry = querry.Where(fmt.Sprintf("%s = ?", _SelfResourceField), userId)
		}

		if res := querry.Delete(_DataEntity); res.RowsAffected == 0 || res.Error != nil {
			handlers.BadRequest(c, "Unable to Delete")
			return
		}

		handlers.StatusOK(c, nil, _Message)
	}
}
