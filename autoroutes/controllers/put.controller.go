package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"gorm.io/gorm"
)

func PutController(db *gorm.DB, _DataEntity interface{}, config models.Config) gin.HandlerFunc {
	// defaultMessageValue := fmt.Sprintf("PUT %s", helpers.GetStructName(_DataEntity))
	// _Message := helpers.ParseArgs(args, "Message", defaultMessageValue).(string)
	// _SelfResource := helpers.ParseArgs(args, "SelfResource", false).(bool)
	// _SelfResourceField := helpers.ParseArgs(args, "SelfResourceField", "user_id").(string)
	// _SelectFields := helpers.ParseArgs(args, "SelectFields", []string{}).([]string)
	// _OmitFields := helpers.ParseArgs(args, "OmitFields", []string{}).([]string)
	// _UpdatableFields := helpers.ParseArgs(args, "UpdatableFields", []string{}).([]string)

	return func(c *gin.Context) {
		entityId := c.Param("id")
		body := helpers.ParseBodyNonStrict(c, config.PutUpdatableFields...)
		if body == nil {
			return
		}

		returningClause := helpers.ReturningColumnsCalculator(db, _DataEntity, config)

		querry := db.Model(_DataEntity).Clauses(returningClause).Where("id  = ?", entityId)

		if config.SelfResource {
			userId, _ := c.Get("user_id")
			querry = querry.Where(fmt.Sprintf("%s = ?", helpers.ToSnake(config.SelfResourceField)), userId)
		}

		if res := querry.Updates(body); res.RowsAffected == 0 || res.Error != nil {
			helpers.BadRequest(c, "Unable to Update / invalid id")
			return
		}

		helpers.StatusOK(c, _DataEntity, config.PutMessage)
	}
}
