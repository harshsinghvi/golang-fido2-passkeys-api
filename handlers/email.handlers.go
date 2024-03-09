package handlers

import (
	"fmt"
	"log"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"gorm.io/gorm"
)

func SendVerificationMail(db *gorm.DB, verification models.Verification) bool {

	var updateVerification = func() bool {
		// TODO: check if wokring
		if res := db.Save(verification); res.RowsAffected == 0 || res.Error != nil {
			log.Println("Error Saving Emai lMessageID, Error, ", res.Error)
			return false
		}
		return true
	}

	verificationUrl, _ := utils.GenerateVerificationUrl(verification.ID.String(), verification.Code)

	if utils.IsEmailDomainTesting(verification.Email) {
		verification.EmailMessageID = verificationUrl
		return updateVerification()
	}

	var toEmail, emailSubject, bodyHtml, bodyHtmlTemplate string

	// TODO USE CONFIG TO STORE HTML BODY TEMPLATE
	switch verification.Type {
	case models.VerificationTypeNewUser:
		emailSubject = "User Verification FIDO 2"
		bodyHtmlTemplate = "<h2> Your User Verification URL :  </h2> <a href=\"%s\">%s</a>"

	case models.VerificationTypeNewPasskey:
		emailSubject = "[Alert] New Passkey registration request"
		bodyHtmlTemplate = "<h2> Your Passkey Authorisation URL :  </h2> <a href=\"%s\">%s</a> <br> please do not authorize this request if yout have not requested this action."

	case models.VerificationTypeDeleteUser:
		emailSubject = "[Alert] User and Data Deletion Request"
		bodyHtmlTemplate = "<h2> Your User and Data Deletion Request Authorisation URL :  </h2> <a href=\"%s\">%s</a> <br> please do not authorize this request if yout have not requested this action."
	default:
		emailSubject = "[Alert/Verification/Authorization] Request from FIDO 2 System"
		bodyHtmlTemplate = "<h2> YourAlert/Verification/Authorization URL :  </h2> <a href=\"%s\">%s</a> <br> please do not authorize this request if yout have not requested this action."
	}

	toEmail = verification.Email
	bodyHtml = fmt.Sprintf(bodyHtmlTemplate, verificationUrl, verificationUrl)

	ok, messageId := utils.SendMailByElasticMail(toEmail, emailSubject, bodyHtml)

	if !ok {
		return false
	}

	verification.EmailMessageID = messageId

	// TODO: check if wokring
	if res := db.Save(verification); res.RowsAffected == 0 || res.Error != nil {
		log.Println("Error Saving EmailMessageID, Error, ", res.Error)
		return false
	}

	return updateVerification()
}
