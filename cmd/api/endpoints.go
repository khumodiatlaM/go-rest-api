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
	router.GET("/users/:id", handlers.AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			userHandler.GetUser(w, r)
		},
		userHandler.JwtSecret,
		userHandler.Logger,
	))

	// ... login user endpoint
	router.POST("/users/login", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		userHandler.LoginUser(w, r)
	})

	return router
}
