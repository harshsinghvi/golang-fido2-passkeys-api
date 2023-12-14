package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"gorm.io/gorm"
)

func PutController(db *gorm.DB, _DataEntity interface{}, config models.Config) gin.HandlerFunc {
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
