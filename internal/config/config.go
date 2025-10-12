package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	Auth0Domain string
	Auth0Aud    string
	Auth0Issuer string
	GeminiKey   string
	GeminiURL   string
	Env         string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("APP_PORT", "8080"),
		DatabaseURL: mustGetEnv("DATABASE_URL"),
		Auth0Domain: mustGetEnv("AUTH0_DOMAIN"),
		Auth0Aud:    mustGetEnv("AUTH0_AUDIENCE"),
		Auth0Issuer: mustGetEnv("AUTH0_ISSUER"),
		GeminiKey:   mustGetEnv("GEMINI_API_KEY"),
		GeminiURL:   mustGetEnv("GEMINI_BASE_URL"),
		Env:         getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("❌ required environment variable %s not set", key)
	}
	return val
}
