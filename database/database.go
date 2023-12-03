package database

import (
	"fmt"
	"log"

	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	var err error
	DB_HOST := utils.GetEnv("DB_HOST", "localhost")
	DB_PORT := utils.GetEnv("DB_PORT", "5432")
	DB_USER := utils.GetEnv("DB_USER", "postgres")
	DB_PASSWORD := utils.GetEnv("DB_PASSWORD", "postgres")
	DB_NAME := utils.GetEnv("DB_NAME", "postgres")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Passkey{})
	DB.AutoMigrate(&models.Challenge{})
	DB.AutoMigrate(&models.AccessToken{})
	// DB.AutoMigrate(&models.PasskeyPrivateKey{})

	tx := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if tx.Error != nil {
		log.Printf("Error in installing PG Extension %s", tx.Error)
	}
	log.Printf("Databse connected and initialised")
}
