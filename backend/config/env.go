package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBPort                 string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	OllamaHost             string
}

var Envs Config = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:              getEnv("DB_ADDRESS", "127.0.0.1"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBName:                 getEnv("DB_NAME", "evoai"),
		JWTExpirationInSeconds: getEnvInteger("JWT_EXP", 60*60*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "not-a-secret"),
		OllamaHost:             getEnv("OLLAMA_HOST", "127.0.0.1:11434"),
	}
}

func getEnv(key string, fallback string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInteger(key string, fallback int64) int64 {

	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}
