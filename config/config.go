package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/TejasThombare20/post-comments-service/utils"
	_ "github.com/lib/pq"
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetDBConfig returns database configuration from environment variables
func GetDBConfig() *DBConfig {
	return &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "post_comments_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// InitDB initializes database connection
func InitDB() (*sql.DB, error) {
	config := GetDBConfig()

	// Build connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	fmt.Println("dsn", dsn)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	utils.LogInfo("Database connected successfully", utils.LogFields{
		"host":   config.Host,
		"port":   config.Port,
		"dbname": config.DBName,
	})
	return db, nil
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
