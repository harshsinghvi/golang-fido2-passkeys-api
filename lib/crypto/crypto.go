package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// Function to generate RSA Public-Private key pair and print in string
func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// Function to encrypt message using public key
func EncryptMessage(message string, publicKey *rsa.PublicKey) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// Function to decrypt ciphertext using private key
func DecryptCipherText(ciphertext []byte, privateKey *rsa.PrivateKey) (string, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// Function to sign message using private key
func SignMessage(message string, privateKey *rsa.PrivateKey) ([]byte, error) {
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, []byte(message))
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Function to verify signature using public key
func VerifySignature(signature []byte, publicKey *rsa.PublicKey, message string) error {
	err := rsa.VerifyPKCS1v15(publicKey, 0, []byte(message), signature)
	return err
}

// Function to parse public key
func ParsePublicKey(pubKeyStr string) (*rsa.PublicKey, error) {
	// Decode base64-encoded string
	decodedPubKey, err := base64.StdEncoding.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}

	// Parse PEM block
	block, _ := pem.Decode(decodedPubKey)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	// Parse public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to assert that parsed public key is of type RSA")
	}

	return rsaPubKey, nil
}

// Function to parse private key
func ParsePrivateKey(privKeyStr string) (*rsa.PrivateKey, error) {
	// Decode base64-encoded string
	decodedPrivKey, err := base64.StdEncoding.DecodeString(privKeyStr)
	if err != nil {
		return nil, err
	}

	// Parse PEM block
	block, _ := pem.Decode(decodedPrivKey)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// Parse private key
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

// Function to convert *rsa.PublicKey into a string
func PublicKeyToString(publicKey *rsa.PublicKey) (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})
	return base64.StdEncoding.EncodeToString(pubPEM), nil
}

// Function to convert *rsa.PrivateKey into a string
func PrivateKeyToString(privateKey *rsa.PrivateKey) (string, error) {
	privASN1 := x509.MarshalPKCS1PrivateKey(privateKey)

	privPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privASN1})
	if privPEM == nil {
		return "", errors.New("failed to encode private key to PEM format")
	}

	return base64.StdEncoding.EncodeToString(privPEM), nil
}

func ValidatePublicAndPrivateKeys(privateKeyStr string, publicKeyStr string) bool {
	privateKey, err := ParsePrivateKey(privateKeyStr)
	if err != nil {
		return false
	}
	publicKey, err := ParsePublicKey(publicKeyStr)

	if err != nil {
		return false
	}

	message := "Secret Message"

	sign, err := SignMessage(message, privateKey)
	if err != nil {
		return false
	}

	err = VerifySignature(sign, publicKey, message)

	if err != nil {
		return false
	}

	return true
}

// Example usage
// privateKey, publicKey, err := GenerateKeyPair()
// if err != nil {
// 	fmt.Println("Error generating key pair:", err)
// 	return
// }

// message := "Hello, Golang RSA Encryption!"

// // Encrypt using public key
// ciphertext, err := EncryptMessage(message, publicKey)
// if err != nil {
// 	fmt.Println("Error encrypting message:", err)
// 	return
// }

// fmt.Println("Encrypted message:", base64.StdEncoding.EncodeToString(ciphertext))

// // Decrypt using private key
// decryptedMessage, err := DecryptCipherText(ciphertext, privateKey)
// if err != nil {
// 	fmt.Println("Error decrypting ciphertext:", err)
// 	return
// }

// fmt.Println("Decrypted message:", decryptedMessage)

// // Sign message using private key
// signature, err := SignMessage(message, privateKey)
// if err != nil {
// 	fmt.Println("Error signing message:", err)
// 	return
// }

// fmt.Println("Message signature:", base64.StdEncoding.EncodeToString(signature))

// // Verify signature using public key
// err = VerifySignature(signature, publicKey, message)
// if err != nil {
// 	fmt.Println("Signature verification failed:", err)
// 	return
// }

// fmt.Println("Signature verified successfully")

// Example public and private key strings (obtained from previous steps or database)
// publicKeyStr := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlGcqJQ4CvZi9WgQvzTbAurpST/L0tuYit8cSIs8cmJiV7E6w7XDbV6+BZmWzoCBXukOVcxCe20cFq8+5Aq0JUPTv/lG6Xv5tSJ9WpOcmK8XC+1jTSMRNjDpSy3jVgI8HmoBJlYXULkR9yIXqRfIi0jU5FWyUSlJlNYfRyXdcsWEBrPq3Xj4N/miUZbpP2S5aXpG8wnTT2zmMk8Rzo0cl3NXwpmOo/TiBc4Tu3U1BF1XM0rQXZ3+fXp/27snB0IjcCW3r7idR/lL98oMzh3U4CTx9TrBQoMLftXL0oFPKMwXZI+zwTPHfJTiC8G6wMSBv+WdU6JHnKntd7fR+k9QIDAQAB"
// privateKeyStr := "MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCUZyolDgK9mL1aBC/NNsC6ulJP8vS25iK3xxIizxyYmJXsTrDtcNtXr4FmZbOgIFe6Q5VzEJ7bRwWrz7kCrQlQ9O/+Ubpe/m1In1ak5yYrxcL7WNNIxE2MOlLLeNWAjweagEmVhdQuRH3IhepF8iLSNTkVbJRKUmk1h9HJd1yxYQGs+rdiZiEaDRHkTLfZJq6S0+tgyW7JzRk3ZqFDmJsscxYrM3AMTyo8zys99wIdg7p9ah9bC3IX3Dxcfo6RPVxxeHsAAleCeaNbJqhzjfuBd44bGojedpLEb77MP9SH0DVsM9yzsL5QDOZz6gMRsX0SdU2o0zOAzee26k0q9AgMBAAECggEAEB+oiEJlyJJSqTWo17wLcwKNw+jpYZNFtb4U73VtdNZHS5D/d2nKnq/8H6qPjprphVgTQ5kgu0fJpRkR7W6F/dJPxGty7YTKqzOthOTWkA/V0GqT0bcql3ljG2KqrnZs5kEBT4IspuolhYDnXMp/d2gN7/WeWS8nxo99j18BQQlNo2ll0+yTEwSiS1/aBbkmEjhMKhRq4aHtV/4Ppmvcq9eDb/Dgh/ODs8Jz3l5LLZMd22kixP9V2tY9ySJuK63W13//dATJ55bTmkYX39pT2zwzWRLB8l3mXQH7naydPPYiFDPK8YVq7HVtCTYjBlRCQGR0HkFkqRq/fniv2OIQKBgQDwmiLcBpBNWuf2qUEmRl6iMDRYuxKyycWRD4BkAmUINmOc5hs41rfXtQ3bVZ7XqD7HzxUzWsfXWt2U1X3L9mV8SJIZQjnZbwvBfgDDDeb+TicfbF5jx58omDgNR18zIufwMfR+j7zIwOb3kTqM/3IDgHvjGpeSvksI05S7uofpFwKBgQCW4B6LeLWj4q43ykvDz5y/FikltX07mfCcqBrpYHZ+jRqWJpLNUb9VWW5mSs5WlRryP1Zv0p1M7roY2Fg8rs5S4sKUK7Ve05bvvwWu0M+FZ1h4gh/cMFOm8HzdtW8EMs2FZPTyL1Vq/QHEHbqE9L+FfsAKwJLKdsE+ol/gNprRo8tQKBgHtppZGEYsQfg3WNCzEvpliOZiPD3tqXvC3ABcTJOdVdM2FvgYrY1m6j0iV8q4WwRDrQ/ZwDsHLiXmOX+3+5qQINAV+bH5BdMw5nWXUnNJMN43R1LLfD6ekTwRMcG52U13tBR5zHAY3xGDXlQFMNmbG2kVOydplS1vMECdBVAoGAKe9KrWDEj0g0lsACZ8t0+DRCjBM+fDzGT6fIM+ZZs+ClVJUTMxFRKknzPTUvRXDvffnJKYFDCpRYp2f8ejdkeZ5xVeGd3nxLf7cf75IfiH4IOfH6RfZu+VXTEf7H50LyA8K8Anx0PryxG3G5bD7UDBtmTzwrzQ0wJ7ZC/Gr4NkCgYBqMDAmH0+0u4vg+r++oBVL3+w/MuGnZZ9HS2R5V44pU/V87WuG8T9F9ZdZdDyIC6T3uByjLh3D97d1SnHhW+/t4vlr/TupppHGRYbBKdymkvSAPgt3r7MCApC1lxK+3T8YqJFq3HdMzgKKIEhtfNcD69lm3+ejSxi2s1E9boUaEmw=="

// // Parse public key
// parsedPublicKey, err := ParsePublicKey(publicKeyStr)
// if err != nil {
// 	fmt.Println("Error parsing public key:", err)
// 	return
// }

// fmt.Println("Parsed Public Key:", parsedPublicKey)

// // Parse private key
// parsedPrivateKey, err := ParsePrivateKey(privateKeyStr)
// if err != nil {
// 	fmt.Println("Error parsing private key:", err)
// 	return
// }

// fmt.Println("Parsed Private Key:", parsedPrivateKey)

// Convert public key to string
// publicKeyStr, err := PublicKeyToString(publicKey)
// if err != nil {
// 	fmt.Println("Error converting public key to string:", err)
// 	return
// }

// fmt.Println("Public Key String:", publicKeyStr)

// // Convert private key to string
// privateKeyStr, err := PrivateKeyToString(privateKey)
// if err != nil {
// 	fmt.Println("Error converting private key to string:", err)
// 	return
// }

// fmt.Println("Private Key String:", privateKeyStr)
