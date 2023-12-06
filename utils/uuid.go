package utils

import (
	"github.com/google/uuid"
)

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
