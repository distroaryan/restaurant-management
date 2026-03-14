package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// If a panic is called instead of crashing the entire server, we call the
// recovery method and log the error.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// We recovered from a panic
				// Log the error and the stack trace
				slog.Error("panic recovered",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
				)	

				// Send a generic 500 Internal Server Error back to the client
				// NEVER send the raw panic error back to the user, as it might leak the sensitive system info
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H {
					"error": "Internal Server error",
				})
			}
		}()

		// Process the request
		c.Next()
	}
}