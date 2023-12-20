package utils

func CheckBoolAndReturnString(b bool, trueString string, falseString string) (bool, string) {
	if b {
		return b, trueString
	}
	return b, falseString
}
