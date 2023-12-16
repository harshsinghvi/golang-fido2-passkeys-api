package main

import (
	"crypto/rsa"
	"fmt"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
	"log"
	"os"
)

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
	message := "Invalid Usage check https://github.com/harshsinghvi/golang-fido2-passkeys-api for cli commands \n"
	message += "cli decrypt -c challenge-string \n"
	message += "cli sign -m challenge-solution \n"
	message += "cli gen \n"
	message += "cli register -n <User Name> -e <User email> --server-url http://localhost:8080 \n"
	message += "cli register-new-key -e <User email> -d <description> --server-url http://localhost:8080 \n"
	message += "cli login --server-url http://localhost:8080 \n"
	message += "cli get-me \n"
	log.Fatalln(message)
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
	if fileExists(CONFIG_PATH) {
		config := readConfigFromFile(CONFIG_PATH)
		if config.ServerUrl != "" {
			return config.ServerUrl
		}
	}
	return DEFAULT_HOST
}

// INFO: Uncomment when needed
// func getPasskeyId(passkeyId string) string {
// 	if passkeyId != "" {
// 		return passkeyId
// 	}
// 	config := readConfigFromFile(CONFIG_PATH)
// 	if config.PasskeyID == "" {
// 		log.Fatal("Server URL Not Found please specify --server-url. using " + DEFAULT_HOST)
// 	}
// 	return config.PasskeyID
// }
