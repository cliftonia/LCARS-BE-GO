// Package database provides PostgreSQL database connection and management.
package database

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // File source driver
)

// MigrateConfig holds configuration for database migrations
type MigrateConfig struct {
	DatabaseURL    string
	MigrationsPath string
}

// RunMigrations runs all pending database migrations
func RunMigrations(cfg MigrateConfig) error {
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		cfg.DatabaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	slog.Info("Database migrations completed",
		"version", version,
		"dirty", dirty,
	)

	return nil
}

// RunMigrationsWithDB runs migrations using an existing database connection
func RunMigrationsWithDB(dbConfig Config, migrationsPath string) error {
	// Build database URL
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)

	return RunMigrations(MigrateConfig{
		DatabaseURL:    dbURL,
		MigrationsPath: migrationsPath,
	})
}

// MigrateDown rolls back the last migration
func MigrateDown(cfg MigrateConfig) error {
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		cfg.DatabaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	slog.Info("Migration rollback completed",
		"version", version,
		"dirty", dirty,
	)

	return nil
}

// MigrateToVersion migrates to a specific version
func MigrateToVersion(cfg MigrateConfig, version uint) error {
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		cfg.DatabaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Migrate(version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate to version %d: %w", version, err)
	}

	slog.Info("Migrated to version", "version", version)

	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(cfg MigrateConfig) (uint, bool, error) {
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		cfg.DatabaseURL,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	version, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get version: %w", err)
	}

	return version, dirty, nil
}

// CreateMigrator creates a new migrate instance for advanced operations
func CreateMigrator(cfg MigrateConfig) (*migrate.Migrate, error) {
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	return m, nil
}

// CreateMigratorWithDriver creates a migrator using the postgres driver directly
func CreateMigratorWithDriver(_ interface{ Driver() interface{} }, migrationsPath string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(nil, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}
