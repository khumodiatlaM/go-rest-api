package http

import (
	"fmt"
	"go-rest-api/pkg/logger"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func StartServer(port string, router *httprouter.Router, logger logger.CustomLogger) {
	// validate port
	if port == "" {
		port = "8080" // default port
	}
	// validate port is an number
	_, err := strconv.Atoi(port)
	if err != nil {
		logger.Fatal("Port must be a number.")
	}
	logger.Info("Starting http server on port " + port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Failed to start http server", "error", err)
	}
}
