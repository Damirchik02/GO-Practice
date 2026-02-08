package main

import (
	"log"
	"net/http"

	"taskmanager/handlers"
	"taskmanager/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTasksHandler(w, r)
		case http.MethodPost:
			handlers.CreateTaskHandler(w, r)
		case http.MethodPatch:
			handlers.UpdateTaskHandler(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"error": "method not allowed"}`))
		}
	})

	// Compose middleware (logging -> apikey -> handlers)
	handler := middleware.LoggingMiddleware(middleware.APIKeyMiddleware(mux))

	log.Println("Server is running on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
