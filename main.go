package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/distroaryan/restaurant-management/internal/config"
	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/logger"
	"github.com/distroaryan/restaurant-management/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// 1. INITIALISE LOGGER
	logger.InitLogger(cfg.Env)
	slog.Info("Starting application...", slog.String("env", cfg.Env), slog.Int("port", cfg.Port))

	// 2. INITIALIZE MONGODB CONNECTION
	db := database.Connect(cfg.MongoURI)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status": "healthy",
		})
	})

	// GRACEFULL SHUTDOWN
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	go func() {
		slog.Info("Server is listening and serving", slog.Int("port", cfg.Port))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to listen and server", slog.String("error", err.Error()))

			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit 
	slog.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("Error", err.Error()))
	}

	// DISCONNECT DATABASE GRACEFULLY
	if err := db.Close(ctx); err != nil {
		slog.Error("Failed to disconnect from MongoDB",
			slog.String("error", err.Error()),
		)
	}

	slog.Info("Server existing")
}
