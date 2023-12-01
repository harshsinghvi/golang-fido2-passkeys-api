package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// TODO: indexex, unique keys https://stackoverflow.com/questions/63409314/how-do-i-create-unique-constraint-for-multiple-columns, defaults, enums
// TODO: Email validation
// Duplicate constraints, savepoints gorm

type User struct {
	gorm.Model
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name  string
	Email string `gorm:"index:idx_email,unique"`
}

type Passkey struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID     uuid.UUID `gorm:"index"`
	Desciption string
	PublicKey  string `gorm:"index:idx_public,unique"`
}

type PasskeyPrivateKey struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID     uuid.UUID `gorm:"index"`
	PasskeyID  uuid.UUID `gorm:"index"`
	PrivateKey string
}

type Challenge struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	PasskeyID uuid.UUID `gorm:"index"`
	UserID    uuid.UUID `gorm:"index"`
	Operand1  int
	Operand2  int
	Operator  string    // `gorm:"type:enum('+','*')"`
	Status    string    // `gorm:"type:enum('FAILED','SUCCESS','PENDING')"`
	Expiry    time.Time `gorm:"index"`
}

type AccessToken struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      uuid.UUID `gorm:"index:idx_access_token"`
	PasskeyID   uuid.UUID `gorm:"index:idx_access_token"`
	ChallengeID uuid.UUID `gorm:"index:idx_access_token"`
	Token       string    `gorm:"index:idx_access_token"`
	Disabled    bool      `gorm:"index:idx_access_token"`
	Expiry      time.Time `gorm:"index:idx_access_token"`
}

type Users []User
type Passkeys []Passkey
type PasskeyPrivateKeys []PasskeyPrivateKey
type Challenges []Challenge
type AccessTokens []AccessToken
