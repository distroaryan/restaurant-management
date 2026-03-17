package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/distroaryan/restaurant-management/internal/config"
	"github.com/distroaryan/restaurant-management/internal/handler"
	"github.com/distroaryan/restaurant-management/internal/middleware"
	"github.com/distroaryan/restaurant-management/internal/routes"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	port       int
}

func NewServer(cfg *config.Config, h *handler.Handler) *Server {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "healthy",
		})
	})

	routes.RegisterRoutes(router, h)

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: router,
		},
		port: cfg.Port,
	}
}

func (s *Server) Start() error {
	slog.Info("Server is listening and serving", slog.Int("port", s.port))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server")
	return s.httpServer.Shutdown(ctx)
}