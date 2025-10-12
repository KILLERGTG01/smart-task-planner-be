package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	Port           string
	DatabaseURL    string
	Auth0Domain    string
	Auth0Aud       string
	Auth0Issuer    string
	GeminiKey      string
	GeminiURL      string
	Env            string
	AllowedOrigins string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:           os.Getenv("APP_PORT"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		Auth0Domain:    os.Getenv("AUTH0_DOMAIN"),
		Auth0Aud:       os.Getenv("AUTH0_AUDIENCE"),
		Auth0Issuer:    os.Getenv("AUTH0_ISSUER"),
		GeminiKey:      os.Getenv("GEMINI_API_KEY"),
		GeminiURL:      os.Getenv("GEMINI_BASE_URL"),
		Env:            os.Getenv("APP_ENV"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
	}

	if cfg.DatabaseURL == "" || cfg.GeminiKey == "" {
		logger, _ := zap.NewProduction()
		logger.Fatal("missing required environment variables")
	}

	return cfg
}
