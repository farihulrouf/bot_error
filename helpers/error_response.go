package helpers

import (
    "encoding/json"
    "net/http"
)

// SendErrorResponse mengirimkan response error dengan pesan dan status code yang ditentukan
func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}