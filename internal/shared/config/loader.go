package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Server
	Port string
	Env  string

	// JWT
	JWTSecret string
	JWTExpiry time.Duration

	// Redis
	RedisURL string

	// Cloudinary
	CloudinaryURL string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "moodle"),

		// Server
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry: parseDuration(getEnv("JWT_EXPIRY", "24h")),

		// Redis
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),

		// Cloudinary
		CloudinaryURL: getEnv("CLOUDINARY_URL", ""),
	}
}

func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword + "@" +
		c.DBHost + ":" + c.DBPort + "/" + c.DBName
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
