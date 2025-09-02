package main

import (
	"go-rest-api/config"
	"go-rest-api/pkg/logger"
)

func main() {
	// ... initialize logger
	logger, err := logger.NewLogger()
	if err != nil {
		logger.Fatal("Failed to initialize logger", "error", err)
	}
	// ... initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config: %v", err)
	}

	logger.Info("Configuration loaded", "config", cfg)
}
