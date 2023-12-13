package helpers

import (
	"fmt"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes/models"
)

func SetDefaultConfig(dEName string, config *models.Config) {
	if config.GetLimit == 0 {
		config.GetLimit = PAGINATION_DEFAULT_LIMIT
	}
	if config.PostDuplicateMessage == "" {
		config.PostDuplicateMessage = "Duplicate Fields."
	}
	if config.SelfResourceField == "" {
		config.SelfResourceField = "user_id"
	}

	if config.GetMessage == "" {
		config.GetMessage = fmt.Sprintf("GET %s", dEName)
	}
	if config.PostMessage == "" {
		config.PostMessage = fmt.Sprintf("POST %s", dEName)
	}
	if config.PutMessage == "" {
		config.PutMessage = fmt.Sprintf("PUT %s", dEName)
	}
	if config.DeleteMessage == "" {
		config.DeleteMessage = fmt.Sprintf("DELETE %s", dEName)
	}
}
