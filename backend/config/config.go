package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string

	RedisHost     string
	RedisPort     string
	RedisPassword string

	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	// Comma-separated allowed origins, e.g. "http://localhost:5173".
	// Defaults to "*" for local dev; set explicitly in production.
	CORSOrigins string

	// Set true in production (HTTPS) to mark refresh cookie as Secure.
	SecureCookies bool
}

func Load() *Config {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env") // fallback when running from backend/

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET env var is required")
	}

	accessTTL, _ := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	refreshTTL, _ := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))

	return &Config{
		Port:             getEnv("PORT", "8080"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "nexus"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       getEnv("POSTGRES_DB", "nexus_messenger"),
		PostgresSSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		JWTSecret:        secret,
		JWTAccessTTL:     accessTTL,
		JWTRefreshTTL:    refreshTTL,
		CORSOrigins:      getEnv("CORS_ORIGINS", "*"),
		SecureCookies:    getEnv("SECURE_COOKIES", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
