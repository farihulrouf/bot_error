package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wagobot.com/model"
)

func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Get self JID from the device store
	/*deviceStore := client.Store.ID
	if deviceStore == nil {
		http.Error(w, "Client not logged in", http.StatusInternalServerError)
		return
	}

	// Convert the deviceStore ID to a proper JID
	selfJID := types.NewJID(deviceStore.User, types.DefaultUserServer)

	// Get user info for the logged-in JID
	userInfoMap, err := client.GetUserInfo([]types.JID{selfJID})
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Ensure the user info is found
	userInfo, exists := userInfoMap[selfJID]
	if !exists {
		http.Error(w, "User info not found", http.StatusNotFound)
		return
	}

	// Prepare the response
	response := map[string]interface{}{
		"device_logged_in": true,
		"self_jid":         selfJID.String(),
		"user_info":        userInfo,
	}

	// Marshal the response into JSON and send it
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToMarshalResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	*/
}

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

/*
func GetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["device"]

	mutex.Lock()
	client := clients[deviceID]
	mutex.Unlock()

	if client == nil {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	if client.IsConnected() {
		// Mendapatkan nomor WhatsApp
		whoami := client.Store.ID.String()
		fmt.Fprintf(w, `{"status": "Connected", "whatsapp_number": "%s"}`, whoami)
	} else {
		fmt.Fprintf(w, `{"status": "Not connected"}`)
	}
}
*/
