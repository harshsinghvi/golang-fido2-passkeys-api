package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harshsinghvi/golang-fido2-passkeys-api/utils/pagination"
)

const (
	MessageBadRequest                 = "Bad Request"
	MessageBadRequestInsufficientData = "Bad Request insufficient data"
	MessageAlreadyVerified            = "Already Verified."
	MessageVerificationAlreadyFailed  = "Verification already Failed."
	MessageVerificationFailed         = "Verification Failed."
	MessageUserVerificationSuccess    = "User Verified."
	MessagePasskeyVerificationSuccess = "Passkey Authorised."
	MessageInvalidVerificationCode    = "Verification Code Invalid."
	MessageExpiredVerificationCode    = "Verification Code Expired."
	MessageErrorWhileSendingEmail     = "Error while sending Email. please register again."
	MessageInvalidBody                = "Invalid Body"
	MessageInvalidEmailAddress        = "Invalid Email Address Please use valid Email."
	MessageInvalidPublicKey           = "Invalid PublicKey Please use valid PublicKey."
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

func UnauthorisedRequest(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"status":  http.StatusUnauthorized,
		"message": "Unauthorised request token invalid/expired/disabled/insufficient roles/permissions.",
	})
	c.Abort()
}
