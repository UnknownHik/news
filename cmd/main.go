package main

import (
	"log"

	"news-rest-api/internal/app"
	"news-rest-api/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
