package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
)

func PutController(_DataEntity interface{}, args ...models.Args) gin.HandlerFunc {
	defaultMessageValue := fmt.Sprintf("PUT %s", utils.GetStructName(_DataEntity))
	_Message := utils.ParseArgs(args, "Message", defaultMessageValue).(string)
	_SelfResource := utils.ParseArgs(args, "SelfResource", false).(bool)
	_SelfResourceField := utils.ParseArgs(args, "SelfResourceField", "user_id").(string)
	_SelectFields := utils.ParseArgs(args, "SelectFields", []string{}).([]string)
	// TODO: Bug can't omit
	_OmitFields := utils.ParseArgs(args, "OmitFields", []string{}).([]string)

	_UpdatableFields := utils.ParseArgs(args, "UpdatableFields", []string{}).([]string)
	return func(c *gin.Context) {
		entityId := c.Param("id")
		body := handlers.ParseBodyNonStrict(c, _UpdatableFields...)
		if body == nil {
			return
		}

		returningClause := helpers.ReturningColumnsCalculator(database.DB, _DataEntity, _SelectFields, _OmitFields)

		querry := database.DB.Model(_DataEntity).Clauses(returningClause).Where("id  = ?", entityId)

		if _SelfResource {
			userId, _ := c.Get("user_id")
			querry = querry.Where(fmt.Sprintf("%s = ?", _SelfResourceField), userId)
		}

		if res := querry.Updates(body); res.RowsAffected == 0 || res.Error != nil {
			handlers.BadRequest(c, "Unable to Update")
			return
		}

		handlers.StatusOK(c, _DataEntity, _Message)
	}
}
