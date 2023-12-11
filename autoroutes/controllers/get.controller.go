package controllers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/pagination"
)

func GetController(_DataEntity interface{}, args ...models.Args) gin.HandlerFunc {
	_Limit := utils.ParseArgs(args, "Limit", pagination.DEFAULT_LIMIT).(int)
	defaultMessageValue := fmt.Sprintf("GET %s", utils.GetStructName(_DataEntity))
	_Message := utils.ParseArgs(args, "Message", defaultMessageValue).(string)
	_OmitFields := utils.ParseArgs(args, "OmitFields", []string{}).([]string)
	_SelectFields := utils.ParseArgs(args, "SelectFields", []string{}).([]string)
	_SelfResource := utils.ParseArgs(args, "SelfResource", false).(bool)
	_SelfResourceField := utils.ParseArgs(args, "SelfResourceField", "user_id").(string)
	_SearchFields := utils.ParseArgs(args, "SearchFields", []string{}).([]string)
	return func(c *gin.Context) {
		var pageStr = c.Query("page")
		var searchStr = c.Query("search")

		pag := pagination.New(pageStr, _Limit)

		if pag.CurrentPage == -1 {
			handlers.BadRequest(c, "Invalid Page.")
			return
		}

		querry := database.DB.Model(_DataEntity)

		if len(_SelectFields) != 0 {
			querry = querry.Select(_SelectFields)
		}

		if len(_OmitFields) != 0 {
			querry = querry.Omit(_OmitFields...)
		}

		if searchStr != "" {
			likeStr := fmt.Sprintf("%%%s%%", searchStr)
			for _, column := range _SearchFields {
				if strings.Contains(column, "id") {
					if utils.IsUUIDValid(searchStr) {
						querry = querry.Or(fmt.Sprintf("%s = ?", column), searchStr)
					}
				} else {
					querry = querry.Or(fmt.Sprintf("%s like ?", column), likeStr)
				}
			}
		}

		if _SelfResource {
			userId, _ := c.Get("user_id")
			querry = querry.Where(fmt.Sprintf("%s = ?", _SelfResourceField), userId)
		}
		res := querry.Count(&pag.TotalRecords)
		if res.Error != nil {
			handlers.BadRequest(c, handlers.MessageBadRequest)
			return
		}
		pag.Validate()
		querry = querry.Order("created_at DESC").Limit(pag.Limit).Offset(pag.Offset)
		res = querry.Find(_DataEntity)
		if res.Error != nil {
			handlers.BadRequest(c, handlers.MessageBadRequest)
			return
		}

		handlers.StatusOKPag(c, _DataEntity, pag, _Message)
	}
}
