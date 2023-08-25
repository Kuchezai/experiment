package main

import (
	"log"

	"experiment.io/config"
	"experiment.io/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load .env file: %s", err)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
