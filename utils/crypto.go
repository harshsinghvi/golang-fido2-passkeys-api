package utils

import (
	"encoding/base64"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
)

func VerifySignature(publicKeyStr string, signatureStr string, message string) bool {
	publicKey, err := crypto.ParsePublicKey(publicKeyStr)
	if err != nil {
		// log.Println("Error While parsing public key from db", err)
		return false
	}

	signature, err := base64.StdEncoding.DecodeString(signatureStr)
	if err != nil {
		// log.Println("Error Parsing signature :", err)
		return false
	}

	err = crypto.VerifySignature(signature, publicKey, message)
	if err != nil {
		// log.Println("Signature verification failed:", err)
		return false
	}

	return true
}
