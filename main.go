package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"harshsinghvi/golang-fido2-passkeys-api/controllers"
	"harshsinghvi/golang-fido2-passkeys-api/database"
	"harshsinghvi/golang-fido2-passkeys-api/utils"
	"log"
)

func init() {
	var err error
	if err = godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
	}
	database.ConnectDb()
	gin.SetMode(gin.DebugMode)
}

func main() {
	PORT := utils.GetEnv("PORT", "8080")

	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/registration/user", controllers.NewUser)
		api.GET("/login/request-challenge/:passkey", controllers.RequestChallenge)
		api.POST("/login/verify-challenge", controllers.VerifyChallenge)

		// login using email and passkey key id
		// auth routes - register new key , check token, business logic
	}

	r.Run(fmt.Sprintf(":%s", PORT))
}