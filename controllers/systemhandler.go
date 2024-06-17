package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"go.mau.fi/whatsmeow"

	"go.mau.fi/whatsmeow/types"
	"wagobot.com/errors"
	"wagobot.com/helpers"
	"wagobot.com/model"
)

func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Get self JID from the device store
	deviceStore := client.Store.ID
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
	deviceStore := client.Store.ID
	if deviceStore == nil {
		http.Error(w, "Client not logged in", http.StatusInternalServerError)
		return
	}

	// Mengonversi ID pengguna menjadi JID
	selfJID := types.NewJID(deviceStore.User, types.DefaultUserServer)

	// Mendapatkan perangkat pengguna yang terhubung dengan JID saat ini
	deviceJIDs, err := client.GetUserDevices([]types.JID{selfJID})
	if err != nil {
		log.Printf("Error getting user devices: %v", err)
		http.Error(w, "Failed to get user devices", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Check device: %v\n", deviceJIDs)
	// Mempersiapkan data respons
	phoneMap := make(map[string]bool)
	responseData := make([]map[string]interface{}, 0)
	for _, deviceJID := range deviceJIDs {
		// Cek jika phone sudah ada di map
		if phoneMap[deviceJID.User] {
			continue // Skip this device if the phone number is already in the map
		}
		/*
			if !helpers.IsLoggedInByNumber(client, deviceJID.User) {
				sendErrorResponse(w, http.StatusBadRequest, "not ready or not available. Please pairing the device", deviceJID.User)
				return
			}
		*/

		// Mendapatkan informasi pengguna untuk setiap perangkat
		userInfoMap, err := client.GetUserInfo([]types.JID{deviceJID})
		fmt.Println("usermap", deviceJID.User)
		if err != nil {
			log.Printf("Error getting user info: %v", err)
			// Menyertakan perangkat dalam respons dengan informasi terbatas
			deviceData := map[string]interface{}{
				"id":      deviceJID.String(),
				"phone":   deviceJID.User,
				"status":  "ready",
				"process": "getMessage",
				"busy":    true,
				"qrcode":  "",
			}
			// Tambahkan phone ke map
			phoneMap[deviceJID.User] = true
			responseData = append(responseData, deviceData)
			continue // Melanjutkan ke perangkat berikutnya
		}

		userInfo := userInfoMap[deviceJID]
		// Print userInfo untuk debugging
		//fmt.Println("User Info:", userInfo)

		deviceData := map[string]interface{}{
			"id":      deviceJID.String(),
			"phone":   deviceJID.User,
			"status":  "ready",
			"name":    "silver",
			"process": "getMessages",
			"busy":    false,
			"qrcode":  "",
		}
		fmt.Println(userInfo)
		// Tambahkan phone ke map
		phoneMap[deviceJID.User] = true
		responseData = append(responseData, deviceData)
	}

	response := map[string]interface{}{
		"data": responseData,
	}

	// Mengubah respons menjadi JSON dan mengirimkannya
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToMarshalResponse)
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
