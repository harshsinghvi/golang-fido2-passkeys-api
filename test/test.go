package main

import (
	"fmt"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
)

// "log"

// "github.com/harshsinghvi/golang-fido2-passkeys-api/database"
// "github.com/joho/godotenv"
// "github.com/harshsinghvi/golang-fido2-passkeys-api/models"
// "github.com/harshsinghvi/golang-fido2-passkeys-api/utils"

func init() {
	// var err error
	// if err = godotenv.Load(); err != nil {
	// 	log.Printf("Error loading .env file")
	// }
	// database.ConnectDb()
}

func main() {
	x := utils.IsEmailDomainTesting("harsh@localhost")

	fmt.Println(x)
	// utils.SendMail()
	// accessToken := models.AccessToken{
	// 	ID: utils.StrToUUID("e5cc62a7-5b7b-4f73-b772-8f2b8be5b999"),
	// }

	// res := database.DB.Delete(&accessToken)

	// log.Print(accessToken)
	// log.Print(res)
	// log.Print(res.RowsAffected)
}
