// Package main is the entry point for the Subspace Backend API server.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/auth"
	"github.com/cliftonbaggerman/subspace-backend/internal/config"
	"github.com/cliftonbaggerman/subspace-backend/internal/constants"
	"github.com/cliftonbaggerman/subspace-backend/internal/database"
	httpserver "github.com/cliftonbaggerman/subspace-backend/internal/http"
	"github.com/cliftonbaggerman/subspace-backend/internal/http/handlers"
	"github.com/cliftonbaggerman/subspace-backend/internal/logger"
	"github.com/cliftonbaggerman/subspace-backend/internal/repository/postgres"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Initialize logger
	log := logger.New(cfg.Server.Environment)

	log.Info("Starting Subspace Backend API",
		"version", cfg.API.Version,
		"environment", cfg.Server.Environment,
	)

	// Initialize database connection
	dbConfig := database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.Name,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Error("Failed to close database connection", "error", err)
		}
	}()

	log.Info("Database connection established")

	// Initialize repositories with PostgreSQL
	userRepo := postgres.NewUserRepository(db)
	messageRepo := postgres.NewMessageRepository(db)

	// Initialize JWT manager
	tokenDuration, _ := time.ParseDuration(cfg.Security.JWTExpiration)
	jwtManager := auth.NewJWTManager(cfg.Security.JWTSecret, tokenDuration)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo)
	authHandler := handlers.NewAuthHandler(userRepo, jwtManager)

	// Create router
	router := httpserver.NewRouter(httpserver.RouterConfig{
		UserHandler:    userHandler,
		MessageHandler: messageHandler,
		AuthHandler:    authHandler,
		UserRepo:       userRepo,
		MessageRepo:    messageRepo,
		JWTManager:     jwtManager,
		CORSOrigins:    cfg.CORS.AllowedOrigins,
	})

	// Create server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  constants.ServerReadTimeout,
		WriteTimeout: constants.ServerWriteTimeout,
		IdleTimeout:  constants.ServerIdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server listening", "address", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Server shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		return
	}

	log.Info("Server stopped gracefully")
}
