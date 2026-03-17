package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/distroaryan/restaurant-management/internal/config"
	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/handler"
	"github.com/distroaryan/restaurant-management/internal/logger"
	"github.com/distroaryan/restaurant-management/internal/server"
)

func main() {
	cfg := config.Load()

	// 1. INITIALISE LOGGER
	logger.InitLogger(cfg.Env)
	slog.Info("Starting application...", slog.String("env", cfg.Env), slog.Int("port", cfg.Port))

	// 2. INITIALIZE MONGODB CONNECTION
	db := database.Connect(cfg.MongoURI)

	appHandler := &handler.Handler{} // Note: Repositories and handlers should be wired up properly here if required, but for now we just pass a handler struct based on existing structure
	appServer := server.NewServer(cfg, appHandler)

	go func() {
		if err := appServer.Start(); err != nil {
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

	if err := appServer.Shutdown(ctx); err != nil {
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
