package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var PORT string
var DefaultDNSAddress string

func Init() {

	_ = godotenv.Load()

	PORT = getEnv("PORT", "8080")
	DefaultDNSAddress = getEnv("DefaultDNSAddress", "1.1.1.1:53")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	fmt.Printf("fallback to default %s\n", key)
	return fallback
}
