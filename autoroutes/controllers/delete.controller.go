package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"gorm.io/gorm"
)

func DeleteController(db *gorm.DB, _DataEntity interface{}, config models.Config) gin.HandlerFunc {
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
