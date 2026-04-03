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
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"

	"kaleidoscope/config"
	"kaleidoscope/controllers"
	"kaleidoscope/database"
	"kaleidoscope/middleware"
	"kaleidoscope/models"
	"kaleidoscope/services"
	"kaleidoscope/telemetry"
)

// Server wraps the HTTP server and dependencies
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
	config     *config.Config
	telemetry  *telemetry.Telemetry
}

// NewServer creates a new HTTP server instance
func NewServer(logger *zap.Logger, config *config.Config) *Server {
	if config.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	tel, err := telemetry.InitTelemetry(context.Background(), config, logger)
	if err != nil {
		logger.Fatal("Failed to initialize telemetry", zap.Error(err))
	}

	db, err := database.Init(config)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	if err := db.DB.Exec(`
		DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'uid') THEN
				ALTER TABLE users ADD COLUMN uid text;
			END IF;
		END $$;
	`).Error; err != nil {
		logger.Fatal("Failed to add uid column", zap.Error(err))
	}

	if err := db.DB.Exec(`
		UPDATE users SET uid = gen_random_uuid()::text 
		WHERE uid IS NULL OR uid = ''
	`).Error; err != nil {
		logger.Fatal("Failed to update existing users with uid", zap.Error(err))
	}

	if err := db.DB.Exec(`
		ALTER TABLE users ALTER COLUMN uid SET NOT NULL;
	`).Error; err != nil {
		logger.Fatal("Failed to set uid not null", zap.Error(err))
	}

	if err := db.DB.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_uid ON users(uid);
	`).Error; err != nil {
		logger.Fatal("Failed to create uid unique index", zap.Error(err))
	}

	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	logger.Info("Database migrations completed successfully")

	userService := services.NewUserService(db.DB)

	var rateLimiter *middleware.RateLimiter
	if config.RateLimit.Enabled {
		rateLimiter = middleware.NewRateLimiter(db.Redis, config.RateLimit.RequestsPerMinute)
		logger.Info("Rate limiter enabled", zap.Int("requests_per_minute", config.RateLimit.RequestsPerMinute))
	}

	router := gin.New()

	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	if config.OTEL.Enabled {
		router.Use(otelgin.Middleware(config.OTEL.ServiceName))
		logger.Info("OpenTelemetry Gin middleware enabled")
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.CORS.AllowOrigins,
		AllowMethods:     config.CORS.AllowMethods,
		AllowHeaders:     config.CORS.AllowHeaders,
		AllowCredentials: config.CORS.AllowCredentials,
	}))

	controllers.RegisterRoutes(router, logger, userService, rateLimiter)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: router,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
		config:     config,
		telemetry:  tel,
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
