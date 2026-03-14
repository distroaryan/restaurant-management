package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Creating a custom logger which intercepts and logs every requests
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		// Process the request
		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)

		status := ctx.Writer.Status()

		// Log the request details
		slog.Info("HTTP Request",
			slog.Int("status", status),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", ctx.ClientIP()),
			slog.String("latency", latency.String()),
		)
	}
}