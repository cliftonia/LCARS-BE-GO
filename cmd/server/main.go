package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/config"
	httpserver "github.com/cliftonbaggerman/subspace-backend/internal/http"
	"github.com/cliftonbaggerman/subspace-backend/internal/http/handlers"
	"github.com/cliftonbaggerman/subspace-backend/internal/repository/memory"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Square Enix Backend API v%s", cfg.API.Version)
	log.Printf("Environment: %s", cfg.Server.Environment)

	// Initialize repositories (in-memory for now)
	userRepo := memory.NewUserRepository()
	messageRepo := memory.NewMessageRepository()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo)

	// Create router
	router := httpserver.NewRouter(httpserver.RouterConfig{
		UserHandler:    userHandler,
		MessageHandler: messageHandler,
		CORSOrigins:    cfg.CORS.AllowedOrigins,
	})

	// Create server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Graceful shutdown with timeout
	// Note: In a real application, you'd want to use context.WithTimeout here
	// but keeping it simple for now

	log.Println("Server stopped")
}
