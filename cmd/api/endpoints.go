package main

import (
	"go-rest-api/internal/handlers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SetupRouter(userHandler *handlers.UserHandler) *httprouter.Router {
	router := httprouter.New()

	// ... health check endpoint
	router.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		handlers.HeathCheck(w, r)
	})

	// ... create user endpoint
	router.POST("/users", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		userHandler.CreateUser(w, r)
	})

	// ... get user by ID endpoint
	router.GET("/users/:id", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		userHandler.GetUser(w, r)
	})

	return router
}
