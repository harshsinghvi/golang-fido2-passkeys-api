package helpers

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

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
