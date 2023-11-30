package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// Function to save RSA private key to a file
func SavePrivateKeyToFile(privateKey *rsa.PrivateKey, filePath string) error {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privKeyBytes})

	err := os.WriteFile(filePath, privKeyPEM, 0644)
	if err != nil {
		return fmt.Errorf("error saving private key to file: %v", err)
	}

	return nil
}

// Function to load RSA private key from a file
func LoadPrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error loading private key from file: %v", err)
	}

	block, _ := pem.Decode(privKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	return privKey, nil
}

// Function to save RSA public key to a file
func SavePublicKeyToFile(publicKey *rsa.PublicKey, filePath string) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("error marshaling public key: %v", err)
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes})

	err = os.WriteFile(filePath, pubKeyPEM, 0644)
	if err != nil {
		return fmt.Errorf("error saving public key to file: %v", err)
	}

	return nil
}

// Function to load RSA public key from a file
func LoadPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error loading public key from file: %v", err)
	}

	block, _ := pem.Decode(pubKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed public key is not of type RSA")
	}

	return rsaPubKey, nil
}

// func Example() {
// 	// Example usage
// 	privateKey, publicKey, err := GenerateKeyPair()
// 	if err != nil {
// 		fmt.Println("Error generating key pair:", err)
// 		return
// 	}

// 	// Save private key to file
// 	privateKeyFilePath := "private_key.pem"
// 	err = SavePrivateKeyToFile(privateKey, privateKeyFilePath)
// 	if err != nil {
// 		fmt.Println("Error saving private key:", err)
// 		return
// 	}
// 	fmt.Println("Private key saved to", privateKeyFilePath)

// 	// Save public key to file
// 	publicKeyFilePath := "public_key.pem"
// 	err = SavePublicKeyToFile(publicKey, publicKeyFilePath)
// 	if err != nil {
// 		fmt.Println("Error saving public key:", err)
// 		return
// 	}
// 	fmt.Println("Public key saved to", publicKeyFilePath)

// 	// Load private key from file
// 	_, err = LoadPrivateKeyFromFile(privateKeyFilePath)
// 	if err != nil {
// 		fmt.Println("Error loading private key:", err)
// 		return
// 	}
// 	fmt.Println("Private key loaded from file")

// 	// Load public key from file
// 	_, err = LoadPublicKeyFromFile(publicKeyFilePath)
// 	if err != nil {
// 		fmt.Println("Error loading public key:", err)
// 		return
// 	}
// 	fmt.Println("Public key loaded from file")
// }
