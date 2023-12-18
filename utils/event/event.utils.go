package event

import (
	"log"
	"strings"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"gorm.io/gorm"
)

const (
	NEW_USER              = "NEW_USER"
	DELETE_USER           = "DELETE_USER"
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	MISC_EVENT            = "MISC_EVENT"
)

func PostEvent(db *gorm.DB, eventName string, data ...string) {
	// TODO: INFO: Webhooks to post event to external service for altering.
	event := models.Event{
		EventName: eventName,
		Data:      strings.Join(data, ", "),
	}

	if res := db.Create(&event); res.Error != nil || res.RowsAffected != 0 {
		log.Println("Error Posting Event to DB ", res.Error)
	}
}
