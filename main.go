package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/controllers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/joho/godotenv"
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

	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/registration/user", controllers.NewUser)
		api.POST("/login/verify-challenge", controllers.VerifyChallenge)
		api.GET("/login/request-challenge/:passkey", controllers.RequestChallenge)
		api.GET("/login/request-challenge", controllers.RequestChallengeUsingPublicKey)
		// TODO: auth routes - register new key , check token, business logic
		// TODO Add healthcheck

		protectedRoutes := api.Group("/protected", controllers.AuthMidlweare())
		{
			protectedRoutes.GET("/get-me", controllers.GetMe)
		}
	}

	router.GET("/health", handlers.HealthHandler)
	// TODO: Pending
	// router.GET("/readiness", handlers.HealthHandler)
	router.GET("/", handlers.ExternalRedirect(REPO_URL))

	router.Run(fmt.Sprintf(":%s", PORT))
}
