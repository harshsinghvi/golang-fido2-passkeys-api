package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
)

const LOCAL_HOST_URL = "http://localhost:8080"

func e(err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 {
			msg = append(msg, "")
		} else {
			msg[0] = fmt.Sprintf("Error while %s, Reason: ", msg[0])
		}
		log.Fatal(msg[0], err.Error())
	}
}

func printError() {
	log.Fatalf("Invalid Usage")
}

func getPrivateKeyFromEnv() string {
	publicKeyStr := os.Getenv("PRIVATE_KEY")
	if publicKeyStr == "" {
		log.Fatalf("PRIVATE_KEY env not set")
	}
	return publicKeyStr
}

func getPrivateKey() *rsa.PrivateKey {
	if ok := fileExists(PRIVATE_KEY_PATH); ok {
		private_key, err := crypto.LoadPrivateKeyFromFile(PRIVATE_KEY_PATH)
		e(err)
		return private_key
	}
	log.Printf("WARN: Passkey (%s) not found Searching for PRIVATE_KEY env.\n", PRIVATE_KEY_PATH)
	privateKeyStr := getPrivateKeyFromEnv()
	privateKey, err := crypto.ParsePrivateKey(privateKeyStr)
	e(err)
	return privateKey
}

func getServerURL(url string) string {
	if url != "" {
		return url
	}
	config := readConfigFromFile(CONFIG_PATH)
	if config.ServerUrl == "" {
		log.Println("Server URL Not Found please specify --server-url. using " + LOCAL_HOST_URL)
		return LOCAL_HOST_URL
	}
	return config.ServerUrl
}

func getPasskeyId(passkeyId string) string {
	if passkeyId != "" {
		return passkeyId
	}
	config := readConfigFromFile(CONFIG_PATH)
	if config.PasskeyID == "" {
		log.Fatal("Server URL Not Found please specify --server-url. using " + LOCAL_HOST_URL)
	}
	return config.PasskeyID
}
