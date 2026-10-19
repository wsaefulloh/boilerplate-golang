package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/wsaefulloh/boilerplate-golang/pkg/logger"
)

type Config struct {
	App      AppConfig
	Database DBConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Env  string
	Port string
	Name string
}

type DBConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	SSLMode     string
	AutoMigrate bool
}

type JWTConfig struct {
	Secret               string
	ExpirationHours      int
	RefreshExpirationDays int
}

type CORSConfig struct {
	AllowedOrigins string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists, ignore error if it doesn't
	if err := godotenv.Load(); err != nil {
		if errRoot := godotenv.Load("../../.env"); errRoot != nil {
			logger.Info("No .env file found or error loading, reading system environment variables")
		}
	}

	cfg := &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			Name: getEnv("APP_NAME", "boilerplate-golang"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:        getEnv("DB_NAME", "boilerplate_db"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			AutoMigrate: getEnvAsBool("DB_AUTO_MIGRATE", false),
		},
		JWT: JWTConfig{
			Secret:                getEnv("JWT_SECRET", "supersecretjwtkeychangeinproduction"),
			ExpirationHours:       getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshExpirationDays: getEnvAsInt("JWT_REFRESH_EXPIRATION_DAYS", 7),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return val
}

func getEnvAsBool(key string, fallback bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valueStr)
	if err != nil {
		return fallback
	}
	return val
}
