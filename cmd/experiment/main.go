package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

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


	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Println("start shutting down server gracefully")
	app.Run(ctx, cfg)
}
