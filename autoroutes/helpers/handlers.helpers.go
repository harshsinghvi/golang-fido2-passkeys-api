package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/pagination"
	"net/http"
)

const (
	MessageBadRequest                 = "Bad Request"
	MessageBadRequestInsufficientData = "Bad Request insufficient data"
	MessageInvalidBody                = "Invalid Body"
	MessageError                      = "Message"
)

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": message,
	})
	c.Abort()
}

func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  http.StatusInternalServerError,
		"message": "Internal server error, Something went wrong !",
	})
	c.Abort()
}

func StatusOK(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": message,
		"data":    data,
	})
	c.Abort()
}

func StatusOKPag(c *gin.Context, data interface{}, pag pagination.Pagination, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    message,
		"data":       data,
		"pagination": pag.Validate(),
	})
	c.Abort()
}
