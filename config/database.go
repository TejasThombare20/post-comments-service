package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/TejasThombare20/post-comments-service/utils"
	_ "github.com/lib/pq"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectDatabase establishes a connection to PostgreSQL database
func ConnectDatabase(config *DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		utils.LogError("Failed to open database connection", err, utils.LogFields{
			"host":   config.Host,
			"port":   config.Port,
			"dbname": config.DBName,
		})
		return nil, utils.WrapError(err, "failed to open database connection")
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		utils.LogError("Failed to ping database", err, utils.LogFields{
			"host":   config.Host,
			"port":   config.Port,
			"dbname": config.DBName,
		})
		return nil, utils.WrapError(err, "failed to ping database")
	}

	utils.LogInfo("Successfully connected to database", utils.LogFields{
		"host":   config.Host,
		"port":   config.Port,
		"dbname": config.DBName,
	})

	return db, nil
}

// CloseDatabase closes the database connection
func CloseDatabase(db *sql.DB) error {
	if db != nil {
		if err := db.Close(); err != nil {
			utils.LogError("Failed to close database connection", err, nil)
			return utils.WrapError(err, "failed to close database connection")
		}
		utils.LogInfo("Database connection closed successfully", nil)
	}
	return nil
}
