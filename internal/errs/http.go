package errs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InternalServerError is a custom error response for 500 Internal Server Error
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}