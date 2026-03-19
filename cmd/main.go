package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/distroaryan/restaurant-management/internal/config"
	"github.com/distroaryan/restaurant-management/internal/database"
	"github.com/distroaryan/restaurant-management/internal/handler"
	"github.com/distroaryan/restaurant-management/internal/observability"
	"github.com/distroaryan/restaurant-management/internal/repository"
	"github.com/distroaryan/restaurant-management/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config " + err.Error())
	}

	shutdownTelemetry := observability.InitTelemetry()
	defer func() {
		_ = shutdownTelemetry(context.Background())
	}()

	// Initialize repository and handlers
	db := database.Connect(cfg.MongoURI, cfg.DbName)

	r := repository.NewRepository(db)
	h := handler.NewHandler(r)

	server := server.NewServer(cfg, h)


	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func ()  {
		if err := server.Start(); err != nil {
			slog.Error("Server failed to start")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown")
	}
	stop()
	cancel()
	slog.Info("Server Exited gracefully")
}
