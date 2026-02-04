package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// HTTP server configuration
type ServerConfig struct {
	Port string
	Host string
}

// Database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

// Logging configuration
type LogConfig struct {
	Level string
}

// load configuration from environment variables
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "transfers"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			MaxConns: int32(getEnvAsInt("DB_MAX_CONNS", 25)),
			MinConns: int32(getEnvAsInt("DB_MIN_CONNS", 5)),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	return cfg, nil
}

// PostgreSQL connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
