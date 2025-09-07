package main

import (
	"go-rest-api/internal/handlers"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(userHandler *handlers.UserHandler) *httprouter.Router {
	router := httprouter.New()

	// ... health check endpoint
	healthPath := "/health"
	router.GET(healthPath, handlers.MetricsMiddleware(
		func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			handlers.HeathCheck(w, r)
		},
		healthPath,
		"GET",
	))

	// ... metrics endpoint
	metricPath := "/metrics"
	router.GET(metricPath, handlers.MetricsMiddleware(
		func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			promhttp.Handler().ServeHTTP(w, r)
		},
		metricPath,
		"GET",
	))

	// ... create user endpoint
	createUserPath := "/users"
	router.POST(createUserPath, handlers.MetricsMiddleware(
		func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			userHandler.CreateUser(w, r)
		},
		createUserPath,
		"POST",
	))

	// ... get user by ID endpoint
	getUserPath := "/users/:id"
	router.GET(getUserPath, handlers.MetricsMiddleware(
		handlers.AuthMiddleware(
			func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				userHandler.GetUser(w, r)
			},
			userHandler.JwtSecret,
			userHandler.Logger,
		),
		getUserPath,
		"GET",
	))

	// ... login user endpoint
	loginUserPath := "/users/login"
	router.POST(loginUserPath, handlers.MetricsMiddleware(
		func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			userHandler.LoginUser(w, r)
		},
		loginUserPath,
		"POST",
	))

	return router
}
