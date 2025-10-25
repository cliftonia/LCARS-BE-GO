// Package config provides configuration management for the Subspace Backend API.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
	API      APIConfig
	Security SecurityConfig
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Port        int
	Host        string
	Environment string
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// CORSConfig contains CORS-specific configuration
type CORSConfig struct {
	AllowedOrigins []string
}

// APIConfig contains API-specific configuration
type APIConfig struct {
	Version   string
	RateLimit int
}

// SecurityConfig contains security-related configuration
type SecurityConfig struct {
	JWTSecret     string
	JWTExpiration string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:        getEnvAsInt("PORT", 8080),
			Host:        getEnv("HOST", "localhost"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Name:     getEnv("DB_NAME", "subspace"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
		},
		API: APIConfig{
			Version:   getEnv("API_VERSION", "v1"),
			RateLimit: getEnvAsInt("API_RATE_LIMIT", 100),
		},
		Security: SecurityConfig{
			JWTSecret:     getEnv("JWT_SECRET", "default-secret-change-in-production"),
			JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Server.Port)
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	if c.Security.JWTSecret == "default-secret-change-in-production" && c.Server.Environment == "production" {
		return fmt.Errorf("JWT secret must be set in production")
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
