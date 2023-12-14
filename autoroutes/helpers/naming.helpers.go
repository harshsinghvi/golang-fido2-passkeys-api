package helpers

import (
	"fmt"
	"reflect"
	// "regexp"
	"github.com/iancoleman/strcase"
	"strings"
)

// var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
// var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

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

//	func ToEndpointNameCase(str string) string {
//		snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
//		snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
//		return strings.ToLower(snake)
//	}

func GetStructName(dataEntity interface{}) string {
	t := reflect.TypeOf(dataEntity)
	split := strings.Split(fmt.Sprint(t), ".")
	name := split[len(split)-1]
	return name
}

// Works Well Without this
// strcase.ConfigureAcronym("ID", "ID")
// strcase.ConfigureAcronym("UserID", "UserID")
// strcase.ConfigureAcronym("PasskeyID", "PasskeyID")
// strcase.ConfigureAcronym("ChallengeID", "ChallengeID")
// strcase.ConfigureAcronym("RequestID", "RequestID")
// strcase.ConfigureAcronym("TokenID", "TokenID")

func ToEndpointNameCase(str string) string {
	return strcase.ToKebab(str)
}

func ToSnake(str string) string {
	return strcase.ToSnake(str)
}

func ToCamel(str string) string {
	return strcase.ToCamel(str)
}
