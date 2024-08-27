package controllers

import (
	"encoding/json"
	"net/http"
)

func setResponse(w http.ResponseWriter, statusCode int, data interface{}) {

	status := "success"
	if statusCode != 200 {
		status = "error"
	}

	response := map[string]interface{}{
		"status": status,
		"data": data,
	}

    // Marshal the data into a pretty JSON format
    jsonResponse, err := json.MarshalIndent(response, "", "")
    if err != nil {
        // If marshalling fails, respond with a 500 Internal Server Error
        http.Error(w, "Failed.", http.StatusInternalServerError)
        return
    }

    // Set the content type to application/json
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    // Write the JSON response
    w.Write(jsonResponse)
}