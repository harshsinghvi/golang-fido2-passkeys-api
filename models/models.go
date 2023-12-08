package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TODO: Check for default values for bools

type User struct {
	gorm.Model
	ID       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Email    string         `gorm:"index:idx_email,unique"`
	Roles    pq.StringArray `gorm:"type:text[]"`
	Name     string
	Verified bool // TODO Update code to check for verified passkeys
}

type Passkey struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID     uuid.UUID `gorm:"index"`
	Desciption string
	PublicKey  string `gorm:"index:idx_public,unique"`
	Verified   bool   // TODO Update code to check for verified passkeys
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

type AccessLog struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	TokenID        uuid.UUID `gorm:"index:idx_access_logs"`
	RequestID      uuid.UUID `gorm:"index:idx_access_logs"`
	Path           string    `gorm:"index:idx_access_logs"`
	ClientIP       string
	Method         string    `gorm:"index:idx_access_logs"`
	StatusCode     int       `gorm:"index:idx_access_logs"`
	BillID         uuid.UUID `gorm:"index:idx_access_logs"`
	Billed         bool
	ResponseTime   int64
	ResponseSize   int
	ServerHostname string
}
