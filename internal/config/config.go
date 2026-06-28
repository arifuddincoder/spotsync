package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Dsn       string
	JwtSecret string
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		Dsn:       os.Getenv("DSN"),
		JwtSecret: os.Getenv("JWT_SECRET_KEY"),
	}
}
