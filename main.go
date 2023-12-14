package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/autoroutes"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/controllers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models/roles"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
	}
	database.ConnectDb()
	gin.SetMode(gin.DebugMode)
}

func main() {
	REPO_URL := utils.GetEnv("REPO_URL", "https://github.com/harshsinghvi/golang-fido2-passkeys-api")
	PORT := utils.GetEnv("PORT", "8080")

	router := gin.Default()

	router.GET("/health", handlers.HealthHandler)
	// TODO: Pending
	// router.GET("/readiness", handlers.HealthHandler)
	router.GET("/", handlers.ExternalRedirect(REPO_URL))

	api := router.Group("/api", controllers.LoggerMW())
	{
		api.POST("/registration/user", controllers.NewUser)
		api.GET("/login/request-challenge", controllers.RequestChallenge)
		api.GET("/login/request-challenge/:passkey", controllers.RequestChallenge)
		api.POST("/login/verify-challenge", controllers.VerifyChallenge)
		api.POST("/register/passkey", controllers.RegistereNewPasskey)

		api.GET("/verify/:id", controllers.Verificaion)
		api.GET("/re-verify/u/:email", controllers.ReVerifyUser)
		api.GET("/re-verify/p", controllers.ReVerifyPasskey)

		adminRouter := api.Group("/admin", controllers.ConfigMW(models.Args{"BillingDisable": true}), controllers.AuthMW(roles.SuperAdmin))
		{
			adminRouter.GET("/verify/passkey/:id", controllers.VerifyPasskey)
			adminRouter.GET("/verify/user/:id", controllers.VerifyUser)

			autoGenRouter := adminRouter.Group("/auto")
			{
				autoroutes.GenerateRoutes(database.DB, autoGenRouter, adminAutoRoutes)
			}
		}
		protectedRouter := api.Group("/protected", controllers.AuthMW())
		{
			protectedRouter.GET("/get-me", controllers.ConfigMW(models.Args{"BillingDisable": false}), controllers.GetMe)
			autoroutes.GenerateRoutes(database.DB, protectedRouter, protectedAutoRoutes)
		}
	}

	router.Run(fmt.Sprintf(":%s", PORT))
}
