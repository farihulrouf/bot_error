package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"wagobot.com/model"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service is available")
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	version := "1.0.0" // Ganti dengan versi sistem yang sesuai
	response := model.VersionResponse{Version: version}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDevicesHandler(w http.ResponseWriter, r *http.Request) {
	type ClientInfo struct {
		ID      string `json:"id"`
		Phone   string `json:"phone"`
		Name    string `json:"name"`
		Status  string `json:"status"`
		Process string `json:"process"`
		Busy    bool   `json:"busy"`
		Qrcode  string `json:"qrcode"`
	}

	mutex.Lock()
	defer mutex.Unlock()
	//clients
	var connectedClients []ClientInfo
	//for key := range clients
	for key, client := range clients {
		if client.IsConnected() {
			whoami := client.Store.ID.String()
			clientInfo := ClientInfo{
				ID:      client.Store.ID.String(),
				Phone:   whoami,
				Name:    key,
				Status:  "Connected",
				Process: "getMessage",
				Busy:    true,
				Qrcode:  " ",
			}
			connectedClients = append(connectedClients, clientInfo)
		}
	}

	response := map[string]interface{}{
		"data": connectedClients,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
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
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Set webhook URL
	webhookURL = reqBody.URL

	// Respond with success message
	resp := model.WebhookResponse{Message: fmt.Sprintf("Webhook set to: %s", webhookURL)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
