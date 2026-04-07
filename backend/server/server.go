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
	"kaleidoscope/worker"
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

	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	logger.Info("Database migrations completed successfully")

	// Create Asynq client for task enqueuing
	asynqClient := worker.NewClient(
		fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		config.Redis.Password,
		config.Redis.DB,
	)

	userService := services.NewUserService(db.DB, asynqClient)

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

	controllers.RegisterRoutes(router, logger, userService, rateLimiter, config, db.DB)

	if config.Server.Environment == "production" {
		staticPath := config.Server.StaticFilesPath
		if _, err := os.Stat(staticPath); err == nil {
			router.NoRoute(func(c *gin.Context) {
				path := c.Request.URL.Path
				filePath := staticPath + path
				if _, err := os.Stat(filePath); os.IsNotExist(err) || path == "/" {
					c.File(staticPath + "/index.html")
				} else {
					c.File(filePath)
				}
			})
			logger.Info("Static files serving enabled", zap.String("path", staticPath))
		} else {
			logger.Warn("Static files path not found, skipping static file serving", zap.String("path", staticPath))
		}
	}

	addr := fmt.Sprintf(":%s", config.Server.Port)
	if config.Server.Host != "" {
		addr = fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	}

	httpServer := &http.Server{
		Addr:    addr,
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
	s.logger.Info("Starting HTTP server", zap.String("host", s.config.Server.Host), zap.String("port", s.config.Server.Port), zap.String("environment", s.config.Server.Environment))

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
