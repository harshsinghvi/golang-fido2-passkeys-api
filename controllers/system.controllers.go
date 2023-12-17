package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/database"
)

func HealthHandler(c *gin.Context) {
	if !database.IsDbReady {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "db connection not ready"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "ok"})
}

func ReadinessHandler(c *gin.Context) {
	if ok := database.Ping(); !ok {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"message": "db connection not ready"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "ok"})
}

func ExternalRedirect(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, url)
	}
}
