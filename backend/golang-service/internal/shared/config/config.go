package config

import (
	"fmt"
	"time"
)

// Config holds all application configuration
type Config struct {
	App        AppConfig
	Server     ServerConfig
	Database   DBConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Cloudinary CloudinaryConfig
}

// CloudinaryConfig holds Cloudinary configuration
type CloudinaryConfig struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadPreset string
	Folder       string
	MaxFileSize  int64 // in bytes
}

// AppConfig holds application-level configuration
type AppConfig struct {
	AppName     string
	Environment string // dev, staging, prod
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DBConfig holds PostgreSQL database configuration
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	AccessSecret     string
	RefreshSecret    string
	AccessExpiryMin  int
	RefreshExpiryMin int
}

// DSN returns the PostgreSQL connection string
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Name, c.User, c.Password, c.SSLMode,
	)
}

// Validate validates the entire configuration
func (c *Config) Validate() error {
	if c.Database.Host == "" || c.Database.Name == "" || c.Database.User == "" {
		return fmt.Errorf("invalid DB config: host/name/user required")
	}
	if c.JWT.AccessSecret == "" || c.JWT.RefreshSecret == "" {
		return fmt.Errorf("invalid JWT config: secrets required")
	}
	return nil
}
