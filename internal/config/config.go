package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = "8111"
)

type Config struct {
	Host     string
	Port     string
	MongoURL string
}

func Read() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = defaultHost
	}

	return &Config{
		Host:     host,
		Port:     port,
		MongoURL: os.Getenv("MONGO_URL"),
	}

}
