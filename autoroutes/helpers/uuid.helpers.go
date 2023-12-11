package helpers

import (
	"github.com/google/uuid"
)

var NilUUID = uuid.Nil

func StrToUUID(s string) uuid.UUID {
	return uuid.Must(uuid.Parse(s))
}

func NewUUID() uuid.UUID {
	return uuid.New()
}

func NewUUIDStr() string {
	return NewUUID().String()
}

func UUIDToStr(id uuid.UUID) string {
	return id.String()
}

func IsUUIDValid(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
