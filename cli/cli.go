package main

import (
	"flag"
	"os"
)

func init() {
	ensureDirectoryExists(BASE_PATH)
}

func main() {
	var serverUrl string
	var userEmail string
	subDecrypt := flag.NewFlagSet("decrypt", flag.PanicOnError)
	challenge := subDecrypt.String("c", "", "Challenge")

	subSign := flag.NewFlagSet("sign", flag.PanicOnError)
	message := subSign.String("m", "", "Message")

	subRegister := flag.NewFlagSet("register", flag.PanicOnError)
	userName := subRegister.String("n", "", "User Full Name")
	subRegister.StringVar(&userEmail, "e", "", "User Email")
	subRegister.StringVar(&serverUrl, "server-url", "", "Server URL")

	subLogin := flag.NewFlagSet("login", flag.PanicOnError)
	subLogin.StringVar(&serverUrl, "server-url", "", "Server URL")

	subAddKey := flag.NewFlagSet("add-key", flag.PanicOnError)
	keyDescription := subAddKey.String("d", "Created From CLI", "Key Description")
	subAddKey.StringVar(&userEmail, "e", "", "User Email")
	subAddKey.StringVar(&serverUrl, "server-url", "", "Server URL")

	if len(os.Args) < 2 {
		printError()
		return
	}

	switch os.Args[1] {
	case "gen":
		gen()
	case "get-me":
		getMe()
	case "decrypt":
		subDecrypt.Parse(os.Args[2:])
		cliDecrypt(*challenge)
	case "sign":
		subSign.Parse(os.Args[2:])
		cliSign(*message)
	case "login":
		subLogin.Parse(os.Args[2:])
		login(serverUrl)
	case "logout":
		logout()
	case "register":
		subRegister.Parse(os.Args[2:])
		userReg(*userName, userEmail, serverUrl)
	case "register-new-key":
		subAddKey.Parse(os.Args[2:])
		addKey(userEmail, *keyDescription, serverUrl)
	default:
		printError()
	}
}
