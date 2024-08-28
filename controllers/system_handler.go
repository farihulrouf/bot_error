package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"wagobot.com/base"
	"wagobot.com/model"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Service is available")
	base.SetResponse(w, http.StatusOK, "Service is available")
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	version := "1.0.0" // Ganti dengan versi sistem yang sesuai
	// response := model.VersionResponse{Version: version}	
	base.SetResponse(w, 0, version)
}

// SetWebhookHandler sets the webhook URL
func SetWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Lock webhookURL to ensure thread safety
	mu.Lock()
	defer mu.Unlock()

	// Parse request body
	var reqBody model.WebhookRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		// http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		base.SetResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Set webhook URL
	webhookURL = reqBody.URL

	// Respond with success message
	resp := model.WebhookResponse{Message: fmt.Sprintf("Webhook set to: %s", webhookURL)}

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(resp)
	base.SetResponse(w, http.StatusOK, resp)
}

func sendPayloadToWebhook(payload interface{}, url string) error {
	// Convert the payload object to a JSON string
	payloadBytes, err := json.Marshal(payload)

	//fmt.Println("cek payload", payloadBytes)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	fmt.Println("data payload", string(payloadBytes))

	// Send the JSON string to the webhook
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send payload to webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response from webhook: %s", resp.Status)
	}

	return nil
}
