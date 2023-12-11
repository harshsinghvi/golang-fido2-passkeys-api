package autoroutes

// TODO Isolate all resources and modules
// TODO get gorm db instance from config (pass db in GenerateRoutes -> each controllers)
// Use arg config Struct to make config not

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/controllers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/helpers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
	// WIP: Remove latter
	// "github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	// "github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	// "github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
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
	// config     Config
}

// https://pkg.go.dev/gopkg.in/mcuadros/go-defaults.v1#section-readme
// https://github.com/mcuadros/go-defaults
type Config struct {
	// Args
	// TODO: rename all arg keys with MethodName Prefix
	// GET PUT POST DELETE
	Message           string
	SelfResource      bool
	SelfResourceField string

	// GET PUT POST
	SelectFields []string // snake_case
	OmitFields   []string // snake_case

	// GET
	Limit        int
	SearchFields []string // snake_case

	// PUT
	UpdatableFields []string // CamelCase

	// POST
	DuplicateMessage string
	// Rename omit in post
	OverrideOmit bool
	NewFields    []string // CamelCase
	// Rename GenerateValues
	GenFields models.GenFields // TODO: Isolate this too // CamelCase
}

type Routes []Route

func New(dataEntity interface{}, args models.Args, methods ...string) Route {
	return Route{
		DataEntity: dataEntity,
		Methods:    methods,
		Args:       args,
	}
}

func GenerateRoutes(db *gorm.DB, router *gin.RouterGroup, routes []Route) {
	var info = map[string]interface{}{}

	for _, route := range routes {
		dEName := helpers.GetStructName(route.DataEntity)
		endpointPath := fmt.Sprintf("/%s", helpers.ToEndpointNameCase(dEName))
		endpointPathWithId := fmt.Sprintf("/%s/:id", helpers.ToEndpointNameCase(dEName))

		for _, method := range route.Methods {
			if method == MethodGet {
				router.GET(endpointPath, controllers.GetController(db, route.DataEntity, route.Args))
			}
			if method == MethodPost {
				router.POST(endpointPath, controllers.PostController(db, route.DataEntity, route.Args))
			}
			if method == MethodPut {
				router.PUT(endpointPathWithId, controllers.PutController(db, route.DataEntity, route.Args))
			}
			if method == MethodDelete {
				router.DELETE(endpointPathWithId, controllers.DeleteController(db, route.DataEntity, route.Args))
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
		helpers.StatusOK(c, _info, "Auto Generated Routes")
	}
}

func ValueWraperGenFunc(val interface{}) models.GenFunc {
	return func(args ...interface{}) interface{} {
		return val
	}
}
