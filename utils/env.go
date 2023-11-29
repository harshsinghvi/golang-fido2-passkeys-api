package utils

import "os"

func GetEnv(varNameString string, defaultValue string) string {
	var varValue string
	if varValue = os.Getenv(varNameString); varValue == "" {
		varValue = defaultValue
	}
	return varValue
}
