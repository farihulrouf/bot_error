package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

// SetupCORS sets up CORS middleware with allowed origins, methods, and headers
func SetupCORS(handler http.Handler) http.Handler {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Fatal("FRONTEND_URL not set in .env file")
	}

	return handlers.CORS(
		handlers.AllowedOrigins([]string{frontendURL}), // Use the URL from .env

		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allow Authorization header for JWT
	)(handler)
}
