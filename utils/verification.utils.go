package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/url"

	"github.com/google/uuid"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

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

func GenerateCode() string {
	return EncodeToString(4)
}

func GenerateVerificationUrl(id string, code string) (string, bool) {
	var BACKEND_URL = GetEnv("BACKEND_URL", "https://passkey.harshsinghvi.com")

	verificationUrl, err := url.Parse(BACKEND_URL)

	if err != nil {
		return "", false
	}

	verificationUrl.Path = fmt.Sprintf("/api/verify/%s", id)
	verificationUrl.RawQuery = fmt.Sprintf("code=%s", code)

	return verificationUrl.String(), true
}

// Args Type, UserID, Email
func CreateVerification(entityId uuid.UUID, args ...interface{}) models.Verification {
	verification := models.Verification{
		EntityID: entityId,
		Status:   models.StatusPending,
		Expiry: GenerateVerificationExpiryDate(),
		Code:   GenerateCode(),
	}

	if len(args) >= 1 {
		verification.Type = args[0].(string)
	}

	if len(args) >= 2 {
		verification.UserID = args[1].(uuid.UUID)
	}

	if len(args) >= 3 {
		verification.Email = args[2].(string)
	}

	return verification
}
