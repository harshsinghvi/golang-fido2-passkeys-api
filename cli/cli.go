package main

import (
	"encoding/base64"
	"flag"
	"harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
	"log"
	"os"
)

func getPrivateKeyFromEnv() string {
	publicKeyStr := os.Getenv("PRIVATE_KEY")
	if publicKeyStr == "" {
		log.Fatalf("PRIVATE_KEY env not set")
	}
	return publicKeyStr
}

func printError() {
	log.Fatalf("Invalid Usage")
}

func gen() (string, string) {
	// Generate Keys
	privateKey, publicKey, err := crypto.GenerateKeyPair()
	if err != nil {
		log.Fatal("Error generating key pair:", err)
	}

	// Convert public key to string
	publicKeyStr, err := crypto.PublicKeyToString(publicKey)
	if err != nil {
		log.Fatal("Error converting public key to string:", err)
	}

	log.Println("Public Key String:", publicKeyStr)

	// Convert private key to string
	privateKeyStr, err := crypto.PrivateKeyToString(privateKey)
	if err != nil {
		log.Fatal("Error converting private key to string:", err)
	}

	log.Println("Private Key String:", privateKeyStr)

	return publicKeyStr, privateKeyStr
}

func decrypt(key string, challenge string) string {
	privateKey, err := crypto.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Error Parsing key : %s", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(challenge)

	if err != nil {
		log.Fatalf("Error Parsing challenge : %s", err)
	}

	decryptedMessage, err := crypto.DecryptCipherText(ciphertext, privateKey)
	if err != nil {
		log.Fatalf("Error decrypting ciphertext: %s", err)
	}

	log.Println("Decrypted challenge:", decryptedMessage)
	return decryptedMessage
}

func sign(key string, message string) string {
	privateKey, err := crypto.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("Error Parsing key : %s", err)
	}
	// Sign message using private key
	signature, err := crypto.SignMessage(message, privateKey)
	if err != nil {
		log.Fatal("Error signing message:", err)
	}

	log.Println("Message signature:", base64.StdEncoding.EncodeToString(signature))
	return base64.StdEncoding.EncodeToString(signature)
}

func main() {
	subDecrypt := flag.NewFlagSet("decrypt", flag.PanicOnError)
	challenge := subDecrypt.String("c", "", "Challenge")

	subSign := flag.NewFlagSet("sign", flag.PanicOnError)
	message := subSign.String("m", "", "Message")

	if len(os.Args) < 2 {
		printError()
		return
	}

	switch os.Args[1] {
	case "gen":
		gen()
	case "decrypt":
		subDecrypt.Parse(os.Args[2:])
		decrypt(getPrivateKeyFromEnv(), *challenge)
	case "sign":
		subSign.Parse(os.Args[2:])
		sign(getPrivateKeyFromEnv(), *message)
	default:
		printError()
	}
}
