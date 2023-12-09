package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/jackc/pgerrcode"
	"gorm.io/gorm/clause"
)

func PostController(_DataEntity interface{}, args ...models.Args) gin.HandlerFunc {
	defaultMessageValue := fmt.Sprintf("POST %s", utils.GetStructName(_DataEntity))
	_Message := utils.ParseArgs(args, "Message", defaultMessageValue).(string)
	_SelfResource := utils.ParseArgs(args, "SelfResource", false).(bool)
	_SelfResourceField := utils.ParseArgs(args, "SelfResourceField", "user_id").(string)
	_NewFields := utils.ParseArgs(args, "NewFields", []string{}).([]string)
	_SelectFields := utils.ParseArgs(args, "SelectFields", []string{}).([]string)
	_GenFields := utils.ParseArgs(args, "GenFields", models.GenFields{}).(models.GenFields)
	_DuplicateMessage := utils.ParseArgs(args, "DuplicateMessage", "Duplicate Fields.").(string)

	return func(c *gin.Context) {
		body := map[string]interface{}{}

		// Generate Fields
		for columnName, valueFunc := range _GenFields {
			body[columnName] = valueFunc(c, body)
		}

		// Parse request body
		reqBody := handlers.ParseBodyStrict(c, _NewFields...)
		if reqBody == nil {
			return
		}

		// Combine genetated fields and request body
		for key, value := range reqBody {
			body[key] = value
		}

		if _SelfResource {
			userId, _ := c.Get("user_id")
			body[_SelfResourceField] = userId
		}

		columns := []clause.Column{}
		for _, column := range _SelectFields {
			columns = append(columns, clause.Column{Name: column})
		}

		body["CreatedAt"] = utils.TimeNow()
		body["UpdatedAt"] = utils.TimeNow()

		if res := database.DB.Model(_DataEntity).Clauses(clause.Returning{Columns: columns}).Create(body); res.RowsAffected == 0 || res.Error != nil {
			switch code, _ := utils.PgErrorCodeAndMessage(res.Error); code {
			case pgerrcode.UniqueViolation:
				handlers.BadRequest(c, _DuplicateMessage)
			default:
				log.Printf("Error While Creating in database: %s", res.Error)
				handlers.InternalServerError(c)
			}
			return
		}

		handlers.StatusOK(c, body, _Message)
	}
}
