package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"kaleidoscope/config"
	"kaleidoscope/controllers"
	"kaleidoscope/database"
	"kaleidoscope/middleware"
	"kaleidoscope/models"
	"kaleidoscope/services"
)

// Server wraps the HTTP server and dependencies
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
	config     *config.Config
}

// NewServer creates a new HTTP server instance
func NewServer(logger *zap.Logger, config *config.Config) *Server {
	// Set Gin mode based on environment
	if config.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database connections
	db, err := database.Init(config)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Run automatic database migrations
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}
	logger.Info("Database migrations completed successfully")

	// Create UserService instance
	userService := services.NewUserService(db.DB)

	// Create Gin engine
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.CORS.AllowOrigins,
		AllowMethods:     config.CORS.AllowMethods,
		AllowHeaders:     config.CORS.AllowHeaders,
		AllowCredentials: config.CORS.AllowCredentials,
	}))

	// Initialize controllers and register routes
	controllers.RegisterRoutes(router, logger, userService)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: router,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
		config:     config,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", zap.String("port", s.config.Server.Port), zap.String("environment", s.config.Server.Environment))

	// Start server in a goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server failed to start", zap.Error(err))
		}
	}()

	return nil
}

// Stop gracefully shuts down the HTTP server
func (s *Server) Stop() error {
	s.logger.Info("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP server forced to shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("HTTP server stopped")
	return nil
}

// WaitForShutdown waits for interrupt signal to gracefully shutdown the server
func (s *Server) WaitForShutdown() {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info("Received shutdown signal")

	// Attempt graceful shutdown
	if err := s.Stop(); err != nil {
		s.logger.Error("Failed to shutdown server gracefully", zap.Error(err))
	}
}
