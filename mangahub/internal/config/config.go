package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string
	AppPort   string
	AppSecret string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DatabaseURL string

	JWTSecret      string
	JWTExpiryHours int

	TCPPort  string
	UDPPort  string
	GRPCPort string

	LogLevel string
}

var App *Config

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	_ = godotenv.Load()

	expiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	cfg := &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		AppPort:   getEnv("APP_PORT", "8080"),
		AppSecret: getEnv("APP_SECRET", "change-me-in-production"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "mangahub"),
		DBPassword: getEnv("DB_PASSWORD", "mangahub_password"),
		DBName:     getEnv("DB_NAME", "mangahub_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:      getEnv("JWT_SECRET", "mangahub-jwt-secret-key-32-chars!!"),
		JWTExpiryHours: expiryHours,

		TCPPort:  getEnv("TCP_PORT", "9001"),
		UDPPort:  getEnv("UDP_PORT", "9002"),
		GRPCPort: getEnv("GRPC_PORT", "9003"),

		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Build DATABASE_URL if not explicitly set
	if url := os.Getenv("DATABASE_URL"); url != "" {
		cfg.DatabaseURL = url
	} else {
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
		)
	}

	App = cfg
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
