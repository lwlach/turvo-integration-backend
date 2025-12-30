package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	loadhandler "github.com/lwlach/turvo-integration-backend/internal/handler/load"
	loadservice "github.com/lwlach/turvo-integration-backend/internal/service/load"
	"github.com/lwlach/turvo-integration-backend/internal/turvo"
)

func main() {
	// Initialize Turvo client
	turvoConfig := turvo.Config{
		BaseURL:      getEnv("TURVO_BASE_URL", "https://my-sandbox.turvo.com"),
		APIKey:       getEnv("TURVO_API_KEY", ""), // Legacy: Use ClientName/ClientSecret instead
		ClientName:   getEnv("TURVO_CLIENT_NAME", ""),
		ClientSecret: getEnv("TURVO_CLIENT_SECRET", ""),
		Username:     getEnv("TURVO_USERNAME", ""),
		Password:     getEnv("TURVO_PASSWORD", ""),
	}

	// Validate required authentication credentials
	if turvoConfig.ClientName == "" || turvoConfig.ClientSecret == "" {
		log.Fatal("TURVO_CLIENT_NAME and TURVO_CLIENT_SECRET are required for authentication")
	}
	if turvoConfig.Username == "" || turvoConfig.Password == "" {
		log.Fatal("TURVO_USERNAME and TURVO_PASSWORD are required for authentication")
	}

	turvoClient := turvo.NewClient(turvoConfig)

	// Initialize service layer
	loadService := loadservice.NewService(turvoClient)

	// Initialize handlers
	loadHandler := loadhandler.NewHandler(loadService)

	// Setup chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Heartbeat("/health"))

	// CORS middleware (basic setup - adjust as needed)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		loadHandler.RegisterRoutes(r)
	})

	// Root endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Turvo Integration Backend API"))
	})

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
