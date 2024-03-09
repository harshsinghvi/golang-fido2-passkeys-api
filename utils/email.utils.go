package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
)

func IsEmailDomainTesting(email string) bool {
	emaiSplit := strings.Split(email, "@")
	domain := emaiSplit[len(emaiSplit)-1]
	return IsTestingDomain(domain)
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// TODO: To update
func SendMailByElasticMail(toEmail, emailSubject, bodyHtml string) (bool, string) {
	var data map[string]interface{}

	var API_KEY = GetEnv("ELASTIC_EMAIL_API_KEY", "")
	var FROM_EMAIL = GetEnv("ELASTIC_FORM_EMAIL", "noreply@harshsinghvi.com")
	var FROM_NAME = GetEnv("ELASTIC_FORM_NAME", "FIDO 2 Passkey de")

	if API_KEY == "" {
		log.Println("Elastic Email Api Key not found pelase check env")
		return false, ""
	}

	url, err := url.Parse("https://api.elasticemail.com/v2/email/send")
	if err != nil {
		return false, ""
	}

	querry := url.Query()
	querry.Set("apikey", API_KEY)
	querry.Set("subject", emailSubject)
	querry.Set("from", FROM_EMAIL)
	querry.Set("fromName", FROM_NAME)
	querry.Set("sender", FROM_NAME)
	querry.Set("senderName", FROM_EMAIL)
	querry.Set("to", toEmail)
	querry.Set("bodyHtml", bodyHtml)
	querry.Set("bodyText", "your verification code: 0000")
	querry.Set("isTransactional", "true")
	querry.Set("trackOpens", "true")
	querry.Set("trackClicks", "true")
	url.RawQuery = querry.Encode()

	resp, err := http.Get(url.String())

	if err != nil {
		return false, ""
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, ""
	}

	err = json.Unmarshal(resBody, &data)

	if err != nil {
		return false, ""
	}

	if resp.StatusCode != http.StatusOK {
		return false, ""
	}

	if success, ok := data["success"]; !ok || success == false {
		return false, ""
	}

	// returning Message id
	return true, fmt.Sprint(data["data"].(map[string]interface{})["messageid"])
}
