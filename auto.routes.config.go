package main

import (
	"net/mail"
	"time"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/lib/crypto"
	AppModels "github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"

	. "github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes"
)

var protectedAutoRoutes = Routes{
	Route{
		DataEntity: &[]AppModels.User{},
		Methods:    []string{MethodGet, MethodPut},
		Config: Config{
			SelfResource:       true,
			SelfResourceField:  "ID",
			PutUpdatableFields: []string{"Name"},
			PostValidationFields: ValidationFields{
				"Email": isEmailValid,
			},
		},
	},
	Route{
		DataEntity: &[]AppModels.Passkey{},
		Methods:    []string{MethodGet, MethodPost, MethodDelete},
		Config: Config{
			SelfResource:         true,
			OmitFields:           []string{"PublicKey"},
			GetSearchFields:      []string{"ID", "UserID", "Desciption", "PublicKey"},
			PostNewFields:        []string{"Desciption", "PublicKey"},
			PostDuplicateMessage: "Public Key already in use",
			PostValidationFields: ValidationFields{
				"PublicKey": isPublicKeyValid,
			},
			PostGenerateValues: GenerateFields{
				"Verified": GenerateConstantValue(false),
			},
		},
	},
	Route{
		DataEntity: &[]AppModels.AccessToken{},
		Methods:    []string{MethodGet, MethodPost, MethodDelete, MethodPut},
		Config: Config{
			SelfResource:         true,
			OmitFields:           []string{"Token"},
			GetSearchFields:      []string{"ID", "UserID", "PasskeyID", "Desciption", "Token"},
			PutUpdatableFields:   []string{"Disabled", "Expiry"},
			PostNewFields:        []string{"Desciption", "Expiry"},
			PostDuplicateMessage: "Public Key already in use",
			PostGenerateValues: GenerateFields{
				"Token":  GenFuncRandomToken,
				"Expiry": TimeNowAfterDays(10),
			},
		},
	},
	Route{
		DataEntity: &[]AppModels.AccessLog{},
		Methods:    []string{MethodGet},
		Config: Config{
			SelfResource:    true,
			GetSearchFields: []string{"ID", "PasskeyID", "TokenID", "BillID", "Method", "Path", "StatusCode"},
		},
	},
}

var adminAutoRoutes = Routes{
	Route{
		DataEntity: &[]AppModels.User{},
		Methods:    []string{MethodGet, MethodPost, MethodDelete, MethodPut},
		Config: Config{
			GetSearchFields:      []string{"ID", "Email", "Name"},
			PostNewFields:        []string{"Name", "Email"},
			PostDuplicateMessage: "Email already in use",
			PostValidationFields: ValidationFields{
				"Email":     isEmailValid,
				"PublicKey": isPublicKeyValid,
			},
			PostGenerateValues: GenerateFields{
				"Verified": GenerateConstantValue(false),
			},
			PutUpdatableFields: []string{"Name", "Roles", "Verified"},
		},
	},

	Route{
		DataEntity: &[]AppModels.Passkey{},
		Methods:    []string{MethodGet, MethodPost, MethodDelete, MethodPut},
		Config: Config{
			GetSearchFields:      []string{"ID", "UserID", "Desciption", "PublicKey"},
			PostNewFields:        []string{"UserID", "Desciption", "PublicKey", "Verified"},
			PutUpdatableFields:   []string{"Desciption", "Verified"},
			PostDuplicateMessage: "Public Key already in use",
			PostValidationFields: ValidationFields{
				"PublicKey": isPublicKeyValid,
			},
		},
	},
	Route{
		DataEntity: &[]AppModels.Challenge{},
		Methods:    []string{MethodGet, MethodDelete},
		Config: Config{
			GetSearchFields: []string{"ID", "UserID", "PasskeyID", "Status"},
		},
	},
	Route{
		DataEntity: &[]AppModels.AccessToken{},
		Methods:    []string{MethodGet, MethodPost, MethodPut, MethodDelete},
		Config: Config{
			GetSearchFields:    []string{"ID", "UserID", "PasskeyID", "ChallengeID", "Desciption", "Token"},
			PostNewFields:      []string{"UserID", "Desciption", "Expiry"},
			PutUpdatableFields: []string{"Disabled", "Expiry", "Desciption"},
			PostGenerateValues: GenerateFields{
				"Token":  GenFuncRandomToken,
				"Expiry": TimeNowAfterDays(10),
			},
		},
	},
	Route{
		DataEntity: &[]AppModels.AccessLog{},
		Methods:    []string{MethodGet},
		Config: Config{
			GetSearchFields: []string{"ID", "UserID", "PasskeyID", "TokenID", "BillID", "Method", "Path", "StatusCode"},
		},
	},
	Route{
		DataEntity: &[]AppModels.Verification{},
		Methods:    []string{MethodGet, MethodPost, MethodPut},
		Config: Config{
			GetSearchFields: []string{"ID", "UserID", "PasskeyID", "ChallengeID", "TokenID", "Code", "Email", "Status", "EmailMessageID"},
			PostNewFields:   []string{"UserID", "Expiry", "Email"},
			PostGenerateValues: GenerateFields{
				"Code":   GenFuncGenerateCode,
				"Expiry": TimeNowAfterDays(10),
				"Status": GenerateConstantValue(AppModels.StatusPending),
			},
			PutUpdatableFields: []string{"Status", "Code", "Expiry"},
		},
	},
}

func GenFuncGenerateCode(args ...interface{}) interface{} {
	return utils.GenerateCode()
}

func GenFuncRandomToken(args ...interface{}) interface{} {
	return utils.GenerateToken(utils.NewUUIDStr())
}

func TimeNowAfterDays(days int) GenerateFunction {
	return func(args ...interface{}) interface{} {
		return time.Now().AddDate(0, 0, days)
	}
}

func TimeNow(args ...interface{}) interface{} {
	return time.Now()
}

func isEmailValid(email interface{}) bool {
	_, err := mail.ParseAddress(email.(string))
	return err == nil
}

func isPublicKeyValid(publicKey interface{}) bool {
	_, err := crypto.ParsePublicKey(publicKey.(string))
	return err == nil
}
