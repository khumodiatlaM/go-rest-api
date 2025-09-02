package main

import (
	"go-rest-api/config"
	"go-rest-api/pkg/database"
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

	// ... connect to database
	db, err := database.Connect(database.DbConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
	}, logger)

	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()
}
