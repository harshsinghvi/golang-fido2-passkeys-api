package utils

import (
	"crypto/rand"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"io"
)

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateCode() string {
	return EncodeToString(4)
}

// TODO: Update Latter
func SendMail(verification models.Verification) bool {
	// verification.Email
	// verification.ID
	// verification.Code
	// Email Body template
	return true
}

func GenFuncGenerateCode(args ...interface{}) interface{} {
	return GenerateCode()
}
