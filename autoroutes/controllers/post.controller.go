package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgerrcode"
	// WIP: TO remove
	// "github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	// "github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	// "github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
)

func PostController(db *gorm.DB, _DataEntity interface{}, args ...models.Args) gin.HandlerFunc {
	defaultMessageValue := fmt.Sprintf("POST %s", helpers.GetStructName(_DataEntity))
	_Message := helpers.ParseArgs(args, "Message", defaultMessageValue).(string)
	_SelfResource := helpers.ParseArgs(args, "SelfResource", false).(bool)
	// TODO Use CamelCase here instead of user_id UserID
	_SelfResourceField := helpers.ParseArgs(args, "SelfResourceField", "user_id").(string)
	_NewFields := helpers.ParseArgs(args, "NewFields", []string{}).([]string)
	_OverrideOmit := helpers.ParseArgs(args, "OverrideOmit", false).(bool)
	_SelectFields := helpers.ParseArgs(args, "SelectFields", []string{}).([]string)
	_GenFields := helpers.ParseArgs(args, "GenFields", models.GenFields{}).(models.GenFields)
	_DuplicateMessage := helpers.ParseArgs(args, "DuplicateMessage", "Duplicate Fields.").(string)
	_OmitFields := helpers.ParseArgs(args, "OmitFields", []string{}).([]string)

	return func(c *gin.Context) {
		body := map[string]interface{}{}

		// Generate Fields
		for columnName, valueFunc := range _GenFields {
			body[columnName] = valueFunc(c, body)
		}

		// Parse request body
		reqBody := helpers.ParseBodyStrict(c, _NewFields...)
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

		returningClause := helpers.ReturningColumnsCalculator(db, _DataEntity, _SelectFields, _OmitFields, _OverrideOmit)

		body["CreatedAt"] = helpers.TimeNow()
		body["UpdatedAt"] = helpers.TimeNow()

		if res := db.Model(_DataEntity).Clauses(returningClause).Create(body); res.RowsAffected == 0 || res.Error != nil {
			switch code, _ := helpers.PgErrorCodeAndMessage(res.Error); code {
			case pgerrcode.UniqueViolation:
				helpers.BadRequest(c, _DuplicateMessage)
			default:
				log.Printf("Error While Creating in database: %s", res.Error)
				helpers.InternalServerError(c)
			}
			return
		}

		// Works Well Without this
		// strcase.ConfigureAcronym("ID", "ID")
		// strcase.ConfigureAcronym("UserID", "UserID")
		// strcase.ConfigureAcronym("PasskeyID", "PasskeyID")
		// strcase.ConfigureAcronym("ChallengeID", "ChallengeID")
		// strcase.ConfigureAcronym("RequestID", "RequestID")
		// strcase.ConfigureAcronym("TokenID", "TokenID")

		for key, value := range body {
			if strings.Contains(key, "id") || strings.Contains(key, "ID") {
				log.Println(key, "====>", value)
				if value != nil {
					body[key] = helpers.StrToUUID(value.(string))
				}
			}
			body[strcase.ToCamel(key)] = value
		}

		x := []map[string]interface{}{body}
		jsonData, err := json.Marshal(x)
		if err != nil {
			helpers.BadRequest(c, helpers.MessageBadRequest)
			return
		}
		err = json.Unmarshal(jsonData, _DataEntity)
		if err != nil {
			helpers.BadRequest(c, helpers.MessageBadRequest)
			return
		}

		// TODO: Returning map[string]intereface{} instead of data entitie's model which is outputs fields in small case
		// TODO: USE PascalCase to snake_case and vice versa functions
		// https://pkg.go.dev/github.com/iancoleman/strcase#section-readme
		// https://www.golinuxcloud.com/go-map-to-struct/
		// Dynamically pass data type to to create single instance or array
		helpers.StatusOK(c, _DataEntity, _Message)
	}
}
