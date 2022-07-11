package main

import (
	"log"

	"github.com/cut4cut/avito-test-work/config"
	"github.com/cut4cut/avito-test-work/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
