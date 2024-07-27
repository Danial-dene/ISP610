// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	AuthToken    string
	ServerPort   string
	Domain			 string
	MailgunAPIKey string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DatabaseURL:  os.Getenv("TURSO_DATABASE_URL"),
		AuthToken:    os.Getenv("TURSO_AUTH_TOKEN"),
		ServerPort:   os.Getenv("PORT"),
		Domain:       os.Getenv("DOMAIN"),
		MailgunAPIKey: os.Getenv("MAILGUN_API_KEY"),
	}
}
