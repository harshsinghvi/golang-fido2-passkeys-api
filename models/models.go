package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// TODO: Email validation,
// TODO: Check for default values for bools

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string
	Email    string `gorm:"index:idx_email,unique"`
	Verified bool   // TODO Update code to check for verified passkeys
}

type Passkey struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID     uuid.UUID `gorm:"index"`
	Desciption string
	Verified   bool   // TODO Update code to check for verified passkeys
	PublicKey  string `gorm:"index:idx_public,unique"`
}

// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
// type PasskeyPrivateKey struct {
// 	gorm.Model
// 	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
// 	UserID     uuid.UUID `gorm:"index"`
// 	PasskeyID  uuid.UUID `gorm:"index"`
// 	PrivateKey string
// }

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
type Challenges []Challenge
type AccessTokens []AccessToken

// INFO: PRIVATE KEY: Uncomment if we need to Store Private Keys
// type PasskeyPrivateKeys []PasskeyPrivateKey

type Args map[string]interface{}
