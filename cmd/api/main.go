package main

import (
	"go-rest-api/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	log.Printf("Config loaded: %+v", cfg)
}
