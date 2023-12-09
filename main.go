package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/controllers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/handlers"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/models/roles"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/autoroutes"
	"github.com/joho/godotenv"
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
	// INFO: CONFIGS
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
		api.POST("/login/verify-challenge", controllers.VerifyChallenge)
		api.GET("/login/request-challenge/:passkey", controllers.RequestChallenge)
		api.GET("/login/request-challenge", controllers.RequestChallengeUsingPublicKey)
		api.POST("/register/passkey", controllers.RegistereNewPasskey)
		api.GET("/verify/:id", controllers.Verificaion)
		api.GET("/re-verify/u/:email", controllers.ReVerifyUser)
		api.GET("/re-verify/p", controllers.ReVerifyPasskey)

		// INFO: Experimantal
		adminRouter := api.Group("/admin", controllers.ConfigMW(models.Args{"BillingDisable": true}), controllers.AuthMW(roles.SuperAdmin))
		{
			adminRouter.GET("/verify/passkey/:id", controllers.VerifyPasskey)
			adminRouter.GET("/verify/user/:id", controllers.VerifyUser)

			autoGenRouter := adminRouter.Group("/auto")
			{
				autoroutes.GenerateRoutes(autoGenRouter, adminAutoRoutes)
			}
		}
		protectedRouter := api.Group("/protected", controllers.AuthMW())
		{
			// INFO: USE OF ConfigMW(models.Args{"BillingDisable": false})
			protectedRouter.GET("/get-me", controllers.GetMe)
			autoroutes.GenerateRoutes(protectedRouter, protectedAutoRoutes)
		}
	}

	router.Run(fmt.Sprintf(":%s", PORT))
}
