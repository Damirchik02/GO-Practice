package middleware

import (
	"fmt"
	"net/http"
	"time"
)

const apiKey = "secret12345" // valid API key

// APIKeyMiddleware checks for a valid API key in the header
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-KEY")
		if key != apiKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "unauthorized"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs the HTTP method, path, and timestamp
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now().Format(time.RFC3339)
		fmt.Printf("%s %s %s\n", ts, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

