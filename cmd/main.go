package main

import (
	"AP2Assignment2/inventory-service/config"
	"AP2Assignment2/inventory-service/internal/app"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Printf("error loading config: %v", err)
		return
	}

	app, err := app.New(ctx, cfg)
	if err != nil {
		log.Printf("error creating app: %v", err)
		return
	}

	err = app.Start()
	if err != nil {
		log.Printf("error starting app: %v", err)
		return
	}
}
