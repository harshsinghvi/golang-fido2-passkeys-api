package main

import (
	"encoding/base64"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
	"log"
)

func gen() (string, string) {
	privateKey, publicKey, err := crypto.GenerateKeyPair() // Generate Keys
	e(err)
	publicKeyStr, err := crypto.PublicKeyToString(publicKey) // Convert public key to string
	e(err)
	privateKeyStr, err := crypto.PrivateKeyToString(privateKey) // Convert private key to string
	e(err)
	err = crypto.SavePrivateKeyToFile(privateKey, PRIVATE_KEY_PATH) // Save private key to file
	e(err)
	err = crypto.SavePublicKeyToFile(publicKey, PUBLIC_KEY_PATH) // Save public key to file
	e(err)

	log.Println("Public Key String:", publicKeyStr)
	// log.Println("Private Key String:", privateKeyStr)
	log.Println("Public key saved to", PUBLIC_KEY_PATH)
	log.Println("Private key saved to", PRIVATE_KEY_PATH)

	writeConfigToFile(Config{}, CONFIG_PATH)
	return publicKeyStr, privateKeyStr
}

func decrypt(challenge string) string {
	privateKey := getPrivateKey()
	ciphertext, err := base64.StdEncoding.DecodeString(challenge)
	e(err)
	decryptedMessage, err := crypto.DecryptCipherText(ciphertext, privateKey)
	e(err)

	return decryptedMessage
}

func sign(message string) string {
	privateKey := getPrivateKey()
	signature, err := crypto.SignMessage(message, privateKey) // Sign message using private key
	e(err)
	return base64.StdEncoding.EncodeToString(signature)
}

func cliDecrypt(challenge string) {
	decryptedMessage := decrypt(challenge)
	log.Println("Decrypted challenge:", decryptedMessage)
}

func cliSign(message string) {
	signature := sign(message)
	log.Println("Message signature:", signature)
}
