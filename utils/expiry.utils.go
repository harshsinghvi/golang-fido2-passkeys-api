package utils

import (
	"time"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/config"
)

func GenerateTokenExpiryDate() time.Time {
	return TimeNowAfterDays(config.TokenExpiryDays)
}

func GenerateChallengeExpiryDate() time.Time {
	return TimeNowAfterDays(config.ChallengeExpiryDays)
}

func GenerateVerificationExpiryDate() time.Time {
	return TimeNowAfterDays(config.VerificationExpiryDays)
}

func TimeNowAfterDays(days int) time.Time {
	return time.Now().AddDate(0, 0, days)
}
