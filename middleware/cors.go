package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// SetupCORS sets up CORS middleware with allowed origins, methods, and headers
func SetupCORS(handler http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8081"}), // Adjust based on your frontend URL
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allow Authorization header for JWT
	)(handler)
}
