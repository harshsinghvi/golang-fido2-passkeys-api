package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
)

func SolveChallenge(challenge models.Challenge) (string, bool) {
	var result int
	switch challenge.Operator {
	case "+":
		result = challenge.Operand1 + challenge.Operand2
	case "*":
		result = challenge.Operand1 * challenge.Operand2
	default:
		return "", false
	}

	return strconv.Itoa(result), true
}

func SolveChallengeString(str string) (string, bool) {
	var challenge models.Challenge
	var err error
	arr := strings.Split(str, " ")

	if challenge.Operand1, err = strconv.Atoi(arr[0]); err != nil {
		return "", false
	}

	if challenge.Operand2, err = strconv.Atoi(arr[2]); err != nil {
		return "", false
	}

	challenge.Operator = arr[1]

	return SolveChallenge(challenge)
}

func CreateChallenge(publicKeyStr string) (string, models.Challenge, error) {
	var challenge models.Challenge
	challenge.Operand1 = rand.Intn(11)
	challenge.Operand2 = rand.Intn(11)
	challenge.Operator = []string{"+", "*"}[rand.Intn(2)]
	publicKey, err := crypto.ParsePublicKey(publicKeyStr)

	if err != nil {
		log.Println("error ParsePublicKey, creating challenge,", err)
		return "", models.Challenge{}, err
	}

	message := fmt.Sprintf("%d %s %d", challenge.Operand1, challenge.Operator, challenge.Operand2)
	ciphertext, err := crypto.EncryptMessage(message, publicKey)

	if err != nil {
		log.Println("error EncryptMessage, creating challenge,", err)
		return "", models.Challenge{}, err
	}

	ciphertextStr := base64.StdEncoding.EncodeToString(ciphertext)
	return ciphertextStr, challenge, nil
}
