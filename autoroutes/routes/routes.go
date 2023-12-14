package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/controllers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	"gorm.io/gorm"
)

func GenerateRoutes(db *gorm.DB, router *gin.RouterGroup, routes []models.Route) {
	var info = map[string]interface{}{}

	for _, route := range routes {

		dEName := helpers.GetStructName(route.DataEntity)
		endpointPath := fmt.Sprintf("/%s", helpers.ToEndpointNameCase(dEName))
		endpointPathWithId := fmt.Sprintf("/%s/:id", helpers.ToEndpointNameCase(dEName))

		helpers.SetDefaultConfig(dEName, &route.Config)

		for _, method := range route.Methods {
			switch method {
			case models.MethodGet:
				router.GET(endpointPath, controllers.GetController(db, route.DataEntity, route.Config))

			case models.MethodPost:
				router.POST(endpointPath, controllers.PostController(db, route.DataEntity, route.Config))

			case models.MethodPut:
				router.PUT(endpointPathWithId, controllers.PutController(db, route.DataEntity, route.Config))

			case models.MethodDelete:
				router.DELETE(endpointPathWithId, controllers.DeleteController(db, route.DataEntity, route.Config))
			}
		}

		info[dEName] = route.Methods
	}

	router.GET("/info", helpers.InfoHandler(info))
}
