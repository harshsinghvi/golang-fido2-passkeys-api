package main

import (
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"path/filepath"
)

var (
	HOME_DIR             = utils.GetEnv("HOME", "/")
	BASE_DIR_NAME        = ".FIDO2"
	PRIVATE_KEY_FILENAME = "passkey.pem"
	PUBLIC_KEY_FILENAME  = "passkey.pub"
	PASSKEY_FILENAME     = "passkey"
	CONFIG_FILENAME      = "config.yml"

	BASE_PATH        = filepath.Join(HOME_DIR, BASE_DIR_NAME)
	PRIVATE_KEY_PATH = filepath.Join(BASE_PATH, PRIVATE_KEY_FILENAME)
	PUBLIC_KEY_PATH  = filepath.Join(BASE_PATH, PUBLIC_KEY_FILENAME)
	PASSKEY_PATH     = filepath.Join(BASE_PATH, PASSKEY_FILENAME)
	CONFIG_PATH      = filepath.Join(BASE_PATH, CONFIG_FILENAME)
)

const (
	// PROD_HOST     = "https://fido2-passkey.onrender.com"
	LOCAL_HOST   = "http://localhost:8080"
	PROD_HOST    = "https://passkey.harshsinghvi.com"
	DEFAULT_HOST = PROD_HOST
)
