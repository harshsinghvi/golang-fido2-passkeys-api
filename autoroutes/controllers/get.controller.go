package controllers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

func GetController(db *gorm.DB, _DataEntity interface{}, config models.Config) gin.HandlerFunc {
	// _Limit := helpers.ParseArgs(args, "Limit", helpers.DEFAULT_LIMIT).(int)
	// defaultMessageValue := fmt.Sprintf("GET %s", helpers.GetStructName(_DataEntity))
	// _Message := helpers.ParseArgs(args, "Message", defaultMessageValue).(string)
	// _OmitFields := helpers.ParseArgs(args, "OmitFields", []string{}).([]string)
	// _SelectFields := helpers.ParseArgs(args, "SelectFields", []string{}).([]string)
	// _SelfResource := helpers.ParseArgs(args, "SelfResource", false).(bool)
	// _SelfResourceField := helpers.ParseArgs(args, "SelfResourceField", "user_id").(string)
	// _SearchFields := helpers.ParseArgs(args, "SearchFields", []string{}).([]string)

	return func(c *gin.Context) {
		var pageStr = c.Query("page")
		var searchStr = c.Query("search")

		pag := helpers.NewPagination(pageStr, config.GetLimit)

		if pag.CurrentPage == -1 {
			helpers.BadRequest(c, "Invalid Page.")
			return
		}

		querry := db.Model(_DataEntity)

		if len(config.SelectFields) != 0 {
			fields := []string{}
			for _, v := range config.SelectFields {
				fields = append(fields, strcase.ToSnake(v))
			}
			querry = querry.Select(fields)
		}

		if len(config.OmitFields) != 0 {
			fields := []string{}
			for _, v := range config.OmitFields {
				fields = append(fields, strcase.ToSnake(v))
			}
			querry = querry.Omit(fields...)
		}

		if searchStr != "" {
			likeStr := fmt.Sprintf("%%%s%%", searchStr)
			for _, columnCamelCase := range config.GetSearchFields {
				column := strcase.ToSnake(columnCamelCase)
				if strings.Contains(column, "id") {
					if helpers.IsUUIDValid(searchStr) {
						querry = querry.Or(fmt.Sprintf("%s = ?", column), searchStr)
					}
				} else {
					querry = querry.Or(fmt.Sprintf("%s like ?", column), likeStr)
				}
			}
		}

		if config.SelfResource {
			userId, _ := c.Get("user_id")
			querry = querry.Where(fmt.Sprintf("%s = ?", strcase.ToSnake(config.SelfResourceField)), userId)
		}
		res := querry.Count(&pag.TotalRecords)
		if res.Error != nil {
			helpers.BadRequest(c, helpers.MessageBadRequest)
			return
		}
		pag.Validate()
		querry = querry.Order("created_at DESC").Limit(pag.Limit).Offset(pag.Offset)
		res = querry.Find(_DataEntity)
		if res.Error != nil {
			helpers.BadRequest(c, helpers.MessageBadRequest)
			return
		}

		helpers.StatusOKPag(c, _DataEntity, pag, config.GetMessage)
	}
}
