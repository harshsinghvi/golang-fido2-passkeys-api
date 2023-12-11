package utils

import "net/mail"

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
