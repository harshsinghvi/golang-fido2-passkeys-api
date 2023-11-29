package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// TODO: indexex, unique keys https://stackoverflow.com/questions/63409314/how-do-i-create-unique-constraint-for-multiple-columns, defaults, enums
// TODO: Email validation

type User struct {
	gorm.Model
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name  string
	Email string
}

type Passkey struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID     uuid.UUID
	Desciption string
	PublicKey  string
}

type PasskeyPrivateKey struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID     uuid.UUID
	PasskeyID  uuid.UUID
	PrivateKey string
}

type Challenge struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	PasskeyID uuid.UUID
	UserID    uuid.UUID
	Operand1  int
	Operand2  int
	Operator  string // `gorm:"type:enum('+','*')"`
	Status    string // `gorm:"type:enum('FAILED','SUCCESS','PENDING')"`
	Expiry    time.Time
}

type AccessToken struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      uuid.UUID
	PasskeyID   uuid.UUID
	ChallengeID uuid.UUID
	Token       string
	Disabled    bool
	Expiry      time.Time
}

type Users []User
type Passkeys []Passkey
type PasskeyPrivateKeys []PasskeyPrivateKey
type Challenges []Challenge
type AccessTokens []AccessToken
