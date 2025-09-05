package main

import (
	"go-rest-api/config"
	"go-rest-api/internal/core"
	userRepo "go-rest-api/internal/db"
	"go-rest-api/internal/handlers"
	"go-rest-api/pkg/database"
	httpserver "go-rest-api/pkg/http"
	"go-rest-api/pkg/logger"
)

func main() {
	// ... initialize logger
	logger, err := logger.NewLogger()
	if err != nil {
		logger.Fatal("Failed to initialize logger", "error", err)
	}

	// ... load the configuration
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

	// ... initialize user repository adapter
	userRepository := userRepo.NewUserRepository(db, logger)

	// ... initialize user service
	userService := core.NewUserService(userRepository, logger)

	// ... initialize user handler
	userHandler := handlers.NewUserHandler(userService)

	// ... setup router
	router := SetupRouter(userHandler)

	// ... start the HTTP server
	httpserver.StartServer(cfg.APIPort, router, logger)
}
