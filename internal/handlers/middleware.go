package handlers

import (
	"context"
	"go-rest-api/internal/metrics"
	"go-rest-api/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func AuthMiddleware(next httprouter.Handle, jwtKey string, logger logger.CustomLogger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("Authorization header is missing")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			logger.Error("Invalid token: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			logger.Error("Invalid token claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims["user_id"].(string)

		// ... add userID to context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		r = r.WithContext(ctx)

		// ... call next handler
		next(w, r, ps)
	}
}

func MetricsMiddleware(next httprouter.Handle, path, method string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		start := time.Now()

		// Use a custom response writer to capture the status code
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next(rw, r, p)

		duration := time.Since(start)

		// Record the metrics
		metrics.RequestCount.WithLabelValues(path, method, http.StatusText(rw.status)).Inc()
		metrics.RequestDuration.WithLabelValues(path, method).Observe(duration.Seconds())
	}
}

// responseWriter is a wrapper to capture the HTTP status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
