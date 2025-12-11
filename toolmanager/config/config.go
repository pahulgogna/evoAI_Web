package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var ToolFile string = "./tools/tools.yml"


func Init() {
	_ = godotenv.Load()

	ToolFile = getEnv(ToolFile, "./tools/tools.yml")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	fmt.Printf("fallback to default %s\n", key)
	return fallback
}