package database

import (
	"fmt"
	"time"

	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConnections holds both read and write database connections
type DatabaseConnections struct {
	Write *gorm.DB
	Read  *gorm.DB
}

// NewConnections creates both read and write database connections
func NewConnections(cfg *config.Config) (*DatabaseConnections, error) {
	// Create write connection
	writeDB, err := createConnection(cfg.GetWriteDSN(), cfg.Database.PG.Write, "write")
	if err != nil {
		return nil, fmt.Errorf("failed to create write connection: %w", err)
	}

	// Create read connection
	readDB, err := createConnection(cfg.GetReadDSN(), cfg.Database.PG.Read, "read")
	if err != nil {
		return nil, fmt.Errorf("failed to create read connection: %w", err)
	}

	return &DatabaseConnections{
		Write: writeDB,
		Read:  readDB,
	}, nil
}

// createConnection creates a database connection with the given configuration
func createConnection(dsn string, connCfg config.PGConnectionConfig, connectionType string) (*gorm.DB, error) {
	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Parse connection lifetime
	maxLifetime, err := time.ParseDuration(connCfg.MaxConnectionLifetime)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse max connection lifetime, using default 5m")
		maxLifetime = 5 * time.Minute
	}

	// Configure connection pool based on config
	sqlDB.SetMaxIdleConns(connCfg.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(connCfg.MaxOpenConnection)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	log.Info().
		Str("type", connectionType).
		Str("host", connCfg.Host).
		Str("port", connCfg.Port).
		Str("database", connCfg.Name).
		Int("max_idle", connCfg.MaxIdleConnection).
		Int("max_open", connCfg.MaxOpenConnection).
		Dur("max_lifetime", maxLifetime).
		Msg("Database connection established")

	return db, nil
}

// Close closes both database connections
func (dc *DatabaseConnections) Close() error {
	var errors []error

	// Close write connection
	if writeDB, err := dc.Write.DB(); err == nil {
		if closeErr := writeDB.Close(); closeErr != nil {
			errors = append(errors, fmt.Errorf("failed to close write connection: %w", closeErr))
		}
	}

	// Close read connection
	if readDB, err := dc.Read.DB(); err == nil {
		if closeErr := readDB.Close(); closeErr != nil {
			errors = append(errors, fmt.Errorf("failed to close read connection: %w", closeErr))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing connections: %v", errors)
	}

	log.Info().Msg("Database connections closed")
	return nil
}

// HealthCheck checks the health of both database connections
func (dc *DatabaseConnections) HealthCheck() error {
	// Check write connection
	if writeDB, err := dc.Write.DB(); err == nil {
		if pingErr := writeDB.Ping(); pingErr != nil {
			return fmt.Errorf("write database health check failed: %w", pingErr)
		}
	} else {
		return fmt.Errorf("failed to get write database instance: %w", err)
	}

	// Check read connection
	if readDB, err := dc.Read.DB(); err == nil {
		if pingErr := readDB.Ping(); pingErr != nil {
			return fmt.Errorf("read database health check failed: %w", pingErr)
		}
	} else {
		return fmt.Errorf("failed to get read database instance: %w", err)
	}

	return nil
}
