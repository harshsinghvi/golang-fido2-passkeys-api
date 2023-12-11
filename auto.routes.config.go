package main

import (
	"time"
	
	AppModels "github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
)

// TODO Check configs
var adminAutoRoutes = autoroutes.Routes{
	autoroutes.New(&[]AppModels.User{}, models.Args{
		"SearchFields":     []string{"id", "email", "name"},
		"UpdatableFields":  []string{"Name", "Roles", "Verified"},
		"NewFields":        []string{"Name", "Email", "Verified"},
		"DuplicateMessage": "Email already in use",
	}, autoroutes.MethodGet, autoroutes.MethodPost, autoroutes.MethodDelete, autoroutes.MethodPut),

	autoroutes.New(&[]AppModels.Passkey{}, models.Args{
		"SearchFields":     []string{"id", "user_id", "desciption", "publicKey"},
		"NewFields":        []string{"UserID", "Desciption", "PublicKey", "Verified"},
		"UpdatableFields":  []string{"Desciption", "Verified"},
		"DuplicateMessage": "Public Key already in use",
	}, autoroutes.MethodGet, autoroutes.MethodPost, autoroutes.MethodDelete, autoroutes.MethodPut),

	autoroutes.New(&[]AppModels.Challenge{}, models.Args{
		"SearchFields": []string{"id", "user_id", "passkey_id", "status"},
	}, autoroutes.MethodGet, autoroutes.MethodDelete),

	autoroutes.New(&[]AppModels.AccessToken{}, models.Args{
		"SearchFields":    []string{"id", "user_id", "passkey_id", "challenge_id", "desciption", "token"},
		"NewFields":       []string{"UserID", "Desciption", "Expiry"},
		"UpdatableFields": []string{"Disabled", "Expiry", "Desciption"},
		"GenFields": models.GenFields{
			"Token":  GenerateRandomToken,
			"Expiry": TimeNowAfterDays(10),
		},
	}, autoroutes.MethodGet, autoroutes.MethodPost, autoroutes.MethodDelete, autoroutes.MethodPut),

	autoroutes.New(&[]AppModels.AccessLog{}, models.Args{
		"SearchFields": []string{"id", "user_id", "passkey_id", "token_id", "bill_id", "method", "path", "status_code"},
	}, autoroutes.MethodGet),

	autoroutes.New(&[]AppModels.Verification{}, models.Args{
		"SearchFields":    []string{"id", "user_id", "passkey_id", "challenge_id", "token_id", "code", "email", "status", "email_message_id"},
		"UpdatableFields": []string{"Status", "Code", "Expiry"},
		"NewFields":       []string{"UserID", "Expiry", "Email"},
		"GenFields": models.GenFields{
			"Code":   utils.GenFuncGenerateCode,
			"Expiry": TimeNowAfterDays(10),
			"Status": autoroutes.ValueWraperGenFunc(AppModels.StatusPending),
		},
	}, autoroutes.MethodGet, autoroutes.MethodPut, autoroutes.MethodPost),
}

var protectedAutoRoutes = autoroutes.Routes{
	autoroutes.New(&[]AppModels.User{}, models.Args{
		"SelfResource":      true,
		"SelfResourceField": "id",
		"UpdatableFields":   []string{"Name"},
	}, autoroutes.MethodGet, autoroutes.MethodPut),

	autoroutes.New(&[]AppModels.Passkey{}, models.Args{
		"SelfResource":     true,
		"OmitFields":       []string{"public_key"},
		"SearchFields":     []string{"id", "user_id", "desciption", "publicKey"},
		"NewFields":        []string{"UserID", "Desciption", "PublicKey"},
		"DuplicateMessage": "Public Key already in use",
	}, autoroutes.MethodGet, autoroutes.MethodPost),

	autoroutes.New(&[]AppModels.AccessToken{}, models.Args{
		"SelfResource": true,
		"SearchFields": []string{"id", "user_id", "passkey_id", "desciption", "token"},
		// "SelectFields":    []string{"id", "passkey_id", "user_id", "disabled", "expiry", "created_at", "updated_at", "desciption"},
		"OmitFields":      []string{"token"},
		"OverrideOmit":    true,
		"UpdatableFields": []string{"Disabled", "Expiry"},
		"NewFields":       []string{"Desciption", "Expiry"},
		"GenFields": models.GenFields{
			"Token":  GenerateRandomToken,
			"Expiry": TimeNowAfterDays(10),
		},
	}, autoroutes.MethodGet, autoroutes.MethodPost, autoroutes.MethodDelete, autoroutes.MethodPut),

	autoroutes.New(&[]AppModels.AccessLog{}, models.Args{
		"SelfResource": true,
		"SearchFields": []string{"id", "passkey_id", "token_id", "bill_id", "method", "path", "status_code"},
	}, autoroutes.MethodGet),
}

func GenerateRandomToken(args ...interface{}) interface{} {
	return utils.GenerateToken(utils.NewUUIDStr())
}

func TimeNowAfterDays(days int) models.GenFunc {
	return func(args ...interface{}) interface{} {
		return time.Now().AddDate(0, 0, days)
	}
}
