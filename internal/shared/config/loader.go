package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Database struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisURL struct {
	Host       string
	Port       string
	Password   string
	DB         int
	TLSEnabled bool
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	Port string
	Env  string

	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration

	RedisURL RedisURL

	CloudinaryURL string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "moodle"),

		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AccessTokenDuration:  parseDuration(getEnv("ACCESS_TOKEN_DURATION", "15m")),
		RefreshTokenDuration: parseDuration(getEnv("REFRESH_TOKEN_DURATION", "168h")),

		RedisURL: RedisURL{
			Host:       getEnv("REDIS_HOST", "localhost"),
			Port:       getEnv("REDIS_PORT", "6379"),
			Password:   getEnv("REDIS_PASSWORD", ""),
			DB:         getEnvAsInt("REDIS_DB", 0),
			TLSEnabled: getEnv("REDIS_TLS_ENABLED", "false") == "true",
		},

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

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
