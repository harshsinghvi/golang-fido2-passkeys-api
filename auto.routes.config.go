package main

import (
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"

	. "github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes"
	. "github.com/harshsinghvi/golang-fido2-passkeys-api/models"
)

// TODO Check configs
var adminAutoRoutes = Routes{
	New(&[]User{}, Args{
		"SearchFields":     []string{"id", "email", "name"},
		"UpdatableFields":  []string{"Name", "Roles", "Verified"},
		"NewFields":        []string{"Name", "Email", "Verified"},
		"DuplicateMessage": "Email already in use",
	}, MethodGet, MethodPost, MethodDelete, MethodPut),

	New(&[]Passkey{}, Args{
		"SearchFields":     []string{"id", "user_id", "desciption", "publicKey"},
		"NewFields":        []string{"UserID", "Desciption", "PublicKey", "Verified"},
		"UpdatableFields":  []string{"Desciption", "Verified"},
		"DuplicateMessage": "Public Key already in use",
	}, MethodGet, MethodPost, MethodDelete, MethodPut),

	New(&[]Challenge{}, Args{
		"SearchFields": []string{"id", "user_id", "passkey_id", "status"},
	}, MethodGet, MethodDelete),

	New(&[]AccessToken{}, Args{
		"SearchFields":    []string{"id", "user_id", "passkey_id", "challenge_id", "desciption", "token"},
		"NewFields":       []string{"UserID", "Desciption", "Expiry"},
		"UpdatableFields": []string{"Disabled", "Expiry", "Desciption"},
		"GenFields": GenFields{
			"Token":  utils.GenerateRandomToken,
			"Expiry": utils.TimeNowAfterDays(10),
		},
	}, MethodGet, MethodPost, MethodDelete, MethodPut),

	New(&[]AccessLog{}, Args{
		"SearchFields": []string{"id", "user_id", "passkey_id", "token_id", "bill_id", "method", "path", "status_code"},
	}, MethodGet),

	New(&[]Verification{}, Args{
		"SearchFields":    []string{"id", "user_id", "passkey_id", "challenge_id", "token_id", "code", "email", "status", "email_message_id"},
		"UpdatableFields": []string{"Status", "Code", "Expiry"},
		"NewFields":       []string{"UserID", "Expiry", "Email"},
		"GenFields": GenFields{
			"Code":   utils.GenFuncGenerateCode,
			"Expiry": utils.TimeNowAfterDays(10),
			"Status": ValueWraperGenFunc(StatusPending),
		},
	}, MethodGet, MethodPut, MethodPost),
}

var protectedAutoRoutes = Routes{
	New(&[]User{}, Args{
		"SelfResource":      true,
		"SelfResourceField": "id",
		"UpdatableFields":   []string{"Name"},
	}, MethodGet, MethodPut),

	New(&[]Passkey{}, Args{
		"SelfResource":     true,
		"OmitFields":       []string{"public_key"},
		"SearchFields":     []string{"id", "user_id", "desciption", "publicKey"},
		"NewFields":        []string{"UserID", "Desciption", "PublicKey"},
		"DuplicateMessage": "Public Key already in use",
	}, MethodGet, MethodPost),

	New(&[]AccessToken{}, Args{
		"SelfResource":    true,
		"SearchFields":    []string{"id", "user_id", "passkey_id", "desciption", "token"},
		"SelectFields":    []string{"id", "passkey_id", "disabled", "expiry", "created_at", "updated_at", "desciption", "token"},
		"OmitFields":      []string{"token", "id"},
		"UpdatableFields": []string{"Disabled", "Expiry"},
		"NewFields":       []string{"Desciption", "Expiry"},
		"GenFields": GenFields{
			"Token":  utils.GenerateRandomToken,
			"Expiry": utils.TimeNowAfterDays(10),
		},
		"OverrideOmit": true,
	}, MethodGet, MethodPost, MethodDelete, MethodPut),

	New(&[]AccessLog{}, Args{
		"SelfResource": true,
		"SearchFields": []string{"id", "passkey_id", "token_id", "bill_id", "method", "path", "status_code"},
	}, MethodGet),
}
