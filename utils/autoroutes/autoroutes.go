package autoroutes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/controllers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
)

const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodDelete = "Delete"
)

type Route struct {
	Methods    []string
	DataEntity interface{}
	Args       models.Args

	// Args
	// Message           string
	// SelfResource      bool
	// SelfResourceField string
	// Limit             int
	// SelectFields      []string
	// SearchFields      []string
	// NewFields         []string // CamelCase
	// GenFields         models.GenFields
	// DuplicateMessage  string
	// UpdatableFields   []string // CamelCase

}

type Routes []Route

func New(dataEntity interface{}, args models.Args, methods ...string) Route {
	return Route{
		DataEntity: dataEntity,
		Methods:    methods,
		Args:       args,
	}
}

func GenerateRoutes(router *gin.RouterGroup, routes []Route) {
	var info = map[string]interface{}{}

	for _, route := range routes {
		dEName := utils.GetStructName(route.DataEntity)
		endpointPath := fmt.Sprintf("/%s", utils.ToEndpointNameCase(dEName))
		endpointPathWithId := fmt.Sprintf("/%s/:id", utils.ToEndpointNameCase(dEName))

		for _, method := range route.Methods {
			if method == MethodGet {
				router.GET(endpointPath, controllers.GetController(route.DataEntity, route.Args))
			}
			if method == MethodPost {
				router.POST(endpointPath, controllers.PostController(route.DataEntity, route.Args))
			}
			if method == MethodPut {
				router.PUT(endpointPathWithId, controllers.PutController(route.DataEntity, route.Args))
			}
			if method == MethodDelete {
				router.DELETE(endpointPathWithId, controllers.DeleteController(route.DataEntity, route.Args))
			}
		}

		info[dEName] = route.Methods
	}

	router.GET("/info", infoHandler(info))
}

func infoHandler(info map[string]interface{}) gin.HandlerFunc {
	_info := map[string]interface{}{}
	for key, value := range info {
		_info[key] = value
	}
	return func(c *gin.Context) {
		handlers.StatusOK(c, _info, "Auto Generated Routes")
	}
}

func ValueWraperGenFunc(val interface{}) models.GenFunc {
	return func(args ...interface{}) interface{} {
		return val
	}
}