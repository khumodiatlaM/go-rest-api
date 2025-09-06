package http

import (
	"context"
	"fmt"
	"go-rest-api/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// ... create a channel to listen for interrupt or terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// ... run the server in a goroutine so that it doesn't block the signal listener
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start http server", "error", err)
		}
	}()
	logger.Info("Starting http server on port " + port)

	// ... block until we receive our signal
	<-quit
	logger.Info("Shutting down server...")

	// ... create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ... attempt a graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server shutting down gracefully")
}
