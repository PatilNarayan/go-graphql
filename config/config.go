package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Initialize environment variables from .env file
func LoadEnv() error {
	if err := godotenv.Load("C:/Users/patil/OneDrive/Desktop/graphql/go-graphql/.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
		return err
	}
	return nil
}

// Helper function to get an environment variable or exit if not set
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}
