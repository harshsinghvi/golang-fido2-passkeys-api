package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
)

// Works Well Without this
// strcase.ConfigureAcronym("ID", "ID")
// strcase.ConfigureAcronym("UserID", "UserID")
// strcase.ConfigureAcronym("PasskeyID", "PasskeyID")
// strcase.ConfigureAcronym("ChallengeID", "ChallengeID")
// strcase.ConfigureAcronym("RequestID", "RequestID")
// strcase.ConfigureAcronym("TokenID", "TokenID")

func PostController(db *gorm.DB, _DataEntity interface{}, config models.Config) gin.HandlerFunc {
	// defaultMessageValue := fmt.Sprintf("POST %s", helpers.GetStructName(_DataEntity))
	// _Message := helpers.ParseArgs(args, "Message", defaultMessageValue).(string)
	// _SelfResource := helpers.ParseArgs(args, "SelfResource", false).(bool)
	// // TODO Use CamelCase here instead of user_id UserID
	// _SelfResourceField := helpers.ParseArgs(args, "SelfResourceField", "user_id").(string)
	// _NewFields := helpers.ParseArgs(args, "NewFields", []string{}).([]string)
	// _OverrideOmit := helpers.ParseArgs(args, "OverrideOmit", false).(bool)
	// _SelectFields := helpers.ParseArgs(args, "SelectFields", []string{}).([]string)
	// _GenFields := helpers.ParseArgs(args, "GenFields", models.GenFields{}).(models.GenFields)
	// _DuplicateMessage := helpers.ParseArgs(args, "DuplicateMessage", "Duplicate Fields.").(string)
	// _OmitFields := helpers.ParseArgs(args, "OmitFields", []string{}).([]string)

	return func(c *gin.Context) {
		body := map[string]interface{}{}

		// Generate Fields
		for columnName, valueFunc := range config.PostGenerateValues {
			body[columnName] = valueFunc(c, body)
		}

		// Parse request body
		reqBody := helpers.ParseBodyStrict(c, config.PostNewFields...)
		if reqBody == nil {
			return
		}

		// Combine genetated fields and request body
		for key, value := range reqBody {
			body[key] = value
		}

		// Validation
		for key, validationFunction := range config.PostValidationFields {
			if !validationFunction(body[key]) {
				helpers.BadRequest(c, fmt.Sprintf(helpers.MessageTemplateInvalidValue, key))
				return
			}
		}

		if config.SelfResource {
			userId, _ := c.Get("user_id")
			body[strcase.ToSnake(config.SelfResourceField)] = userId
		}

		returningClause := helpers.ReturningColumnsCalculator(db, _DataEntity, config)
		log.Println("returningClause ======>", returningClause)

		body["CreatedAt"] = helpers.TimeNow()
		body["UpdatedAt"] = helpers.TimeNow()

		if res := db.Model(_DataEntity).Clauses(returningClause).Create(body); res.RowsAffected == 0 || res.Error != nil {
			switch code, _ := helpers.PgErrorCodeAndMessage(res.Error); code {
			case pgerrcode.UniqueViolation:
				helpers.BadRequest(c, config.PostDuplicateMessage)
			default:
				log.Printf("Error While Creating in database: %s", res.Error)
				helpers.InternalServerError(c)
			}
			return
		}

		for key, value := range body {
			if strings.Contains(key, "id") || strings.Contains(key, "ID") {
				if value != nil {
					body[key] = helpers.StrToUUID(value.(string))
				}
			}
			body[strcase.ToCamel(key)] = value
		}

		// https://www.golinuxcloud.com/go-map-to-struct/
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

		helpers.StatusOK(c, _DataEntity, config.PostMessage)
	}
}
