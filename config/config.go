package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TejasThombare20/post-comments-service/utils"
	_ "github.com/lib/pq"
)

// Config holds all application configuration
type Config struct {
	Database *DBConfig
	Server   *ServerConfig
	JWT      *JWTConfig
	App      *AppConfig
}

// DBConfig holds database configuration
type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// AppConfig holds general application configuration
type AppConfig struct {
	Environment string
	LogLevel    string
	Debug       bool
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("config validation error for %s: %s", e.Field, e.Message)
}

// LoadConfig loads and validates all configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Database: loadDBConfig(),
		Server:   loadServerConfig(),
		JWT:      loadJWTConfig(),
		App:      loadAppConfig(),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// loadDBConfig loads database configuration from environment variables
func loadDBConfig() *DBConfig {
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	connMaxLifetime, _ := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m"))

	return &DBConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "post_comments_db"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}
}

// loadServerConfig loads server configuration from environment variables
func loadServerConfig() *ServerConfig {
	readTimeout, _ := time.ParseDuration(getEnv("SERVER_READ_TIMEOUT", "15s"))
	writeTimeout, _ := time.ParseDuration(getEnv("SERVER_WRITE_TIMEOUT", "15s"))
	idleTimeout, _ := time.ParseDuration(getEnv("SERVER_IDLE_TIMEOUT", "60s"))

	return &ServerConfig{
		Port:         getEnv("PORT", "8080"),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

// loadJWTConfig loads JWT configuration from environment variables
func loadJWTConfig() *JWTConfig {
	accessTokenDuration, _ := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_DURATION", "15m"))
	refreshTokenDuration, _ := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_DURATION", "168h"))

	return &JWTConfig{
		SecretKey:            getEnv("JWT_SECRET_KEY", ""),
		AccessTokenDuration:  accessTokenDuration,
		RefreshTokenDuration: refreshTokenDuration,
	}
}

// loadAppConfig loads general application configuration from environment variables
func loadAppConfig() *AppConfig {
	debug, _ := strconv.ParseBool(getEnv("DEBUG", "false"))

	return &AppConfig{
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Debug:       debug,
	}
}

// validateConfig validates all configuration values
func validateConfig(config *Config) error {
	var errors []ValidationError

	// Validate database configuration
	if config.Database.Host == "" {
		errors = append(errors, ValidationError{"DB_HOST", "database host is required"})
	}
	if config.Database.Port == "" {
		errors = append(errors, ValidationError{"DB_PORT", "database port is required"})
	}
	if config.Database.User == "" {
		errors = append(errors, ValidationError{"DB_USER", "database user is required"})
	}
	if config.Database.DBName == "" {
		errors = append(errors, ValidationError{"DB_NAME", "database name is required"})
	}
	if config.Database.MaxOpenConns <= 0 {
		errors = append(errors, ValidationError{"DB_MAX_OPEN_CONNS", "must be greater than 0"})
	}
	if config.Database.MaxIdleConns <= 0 {
		errors = append(errors, ValidationError{"DB_MAX_IDLE_CONNS", "must be greater than 0"})
	}

	// Validate server configuration
	if config.Server.Port == "" {
		errors = append(errors, ValidationError{"PORT", "server port is required"})
	}
	if port, err := strconv.Atoi(config.Server.Port); err != nil || port <= 0 || port > 65535 {
		errors = append(errors, ValidationError{"PORT", "must be a valid port number (1-65535)"})
	}

	// Validate JWT configuration
	if config.JWT.SecretKey == "" {
		errors = append(errors, ValidationError{"JWT_SECRET_KEY", "JWT secret key is required"})
	}
	if len(config.JWT.SecretKey) < 32 {
		errors = append(errors, ValidationError{"JWT_SECRET_KEY", "JWT secret key must be at least 32 characters long"})
	}
	if config.JWT.AccessTokenDuration <= 0 {
		errors = append(errors, ValidationError{"JWT_ACCESS_TOKEN_DURATION", "must be greater than 0"})
	}
	if config.JWT.RefreshTokenDuration <= 0 {
		errors = append(errors, ValidationError{"JWT_REFRESH_TOKEN_DURATION", "must be greater than 0"})
	}

	// Validate app configuration
	validEnvironments := []string{"development", "staging", "production"}
	if !contains(validEnvironments, config.App.Environment) {
		errors = append(errors, ValidationError{"ENVIRONMENT", fmt.Sprintf("must be one of: %s", strings.Join(validEnvironments, ", "))})
	}

	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, config.App.LogLevel) {
		errors = append(errors, ValidationError{"LOG_LEVEL", fmt.Sprintf("must be one of: %s", strings.Join(validLogLevels, ", "))})
	}

	if len(errors) > 0 {
		return &ConfigValidationError{Errors: errors}
	}

	return nil
}

// ConfigValidationError represents multiple validation errors
type ConfigValidationError struct {
	Errors []ValidationError
}

func (e *ConfigValidationError) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return fmt.Sprintf("configuration validation failed:\n  - %s", strings.Join(messages, "\n  - "))
}

// GetDBConfig returns database configuration from environment variables (legacy function)
func GetDBConfig() *DBConfig {
	config, _ := LoadConfig()
	return config.Database
}

// InitDB initializes database connection with enhanced configuration
func InitDB() (*sql.DB, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	dbConfig := config.Database

	// Build connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode)

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
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

	utils.LogInfo("Database connected successfully", utils.LogFields{
		"host":              dbConfig.Host,
		"port":              dbConfig.Port,
		"dbname":            dbConfig.DBName,
		"max_open_conns":    dbConfig.MaxOpenConns,
		"max_idle_conns":    dbConfig.MaxIdleConns,
		"conn_max_lifetime": dbConfig.ConnMaxLifetime,
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

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
