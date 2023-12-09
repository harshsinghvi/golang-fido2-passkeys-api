package autoroutes

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/controllers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/pagination"
)

type Route struct {
	DataEntity   interface{}
	Limit        int
	SearchFields []string
}

type Routes []Routes

var info map[string]string

func NewRoute(dataEntity interface{}, args ...string) Route {
	return Route{
		DataEntity:   dataEntity,
		SearchFields: args,
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// func ToSnakeCase(str string) string {
// 	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
// 	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
// 	return strings.ToLower(snake)
// }

// func ToNormalPruralName(str string) string {
// 	snake := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
// 	snake = matchAllCap.ReplaceAllString(snake, "${1} ${2}")
// 	return strings.ToLower(snake + "s")
// }

func ToEndpointNameCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

func GetStructName(dataEntity interface{}) string {
	t := reflect.TypeOf(dataEntity)
	split := strings.Split(fmt.Sprint(t), ".")
	name := split[len(split)-1]
	return name
}

func GenerateRoutes(router *gin.RouterGroup, routes []Route) {
	info = map[string]string{}

	for _, route := range routes {
		dEName := GetStructName(route.DataEntity)
		endpointPath := fmt.Sprintf("/%s", ToEndpointNameCase(dEName))

		args := models.Args{
			"SearchFields": route.SearchFields,
			"Message":      dEName,
			"Limit":        route.Limit,
		}

		if route.Limit <= 0 {
			args["Limit"] = pagination.DEFAULT_LIMIT
		}

		router.GET(endpointPath, controllers.GetController(route.DataEntity, args))
		info[endpointPath] = dEName
	}

	router.GET("/info", func(c *gin.Context) {
		handlers.StatusOK(c, info, "Auto Generated Routes")
	})
}
