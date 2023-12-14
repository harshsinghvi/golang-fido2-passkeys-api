package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"gorm.io/gorm"
)

func DeleteController(db *gorm.DB, _DataEntity interface{}, config models.Config) gin.HandlerFunc {
	// defaultMessageValue := fmt.Sprintf("DELETE %s", helpers.GetStructName(_DataEntity))
	// _Message := helpers.ParseArgs(args, "Message", defaultMessageValue).(string)
	// _SelfResource := helpers.ParseArgs(args, "SelfResource", false).(bool)
	// _SelfResourceField := helpers.ParseArgs(args, "SelfResourceField", "user_id").(string)

	return func(c *gin.Context) {
		entityId := c.Param("id")

		querry := db.Where("id = ?", entityId)

		if config.SelfResource {
			userId, _ := c.Get("user_id")
			querry = querry.Where(fmt.Sprintf("%s = ?", helpers.ToSnake(config.SelfResourceField)), userId)
		}

		if res := querry.Delete(_DataEntity); res.RowsAffected == 0 || res.Error != nil {
			helpers.BadRequest(c, "Unable to Delete / invalid id")
			return
		}

		helpers.StatusOK(c, nil, config.DeleteMessage)
	}
}
