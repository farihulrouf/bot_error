package controllers

import (
	// "encoding/json"
	"fmt"
	"net/http"
	// "context"
	// "log"
	// "time"
)

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

	var connectedClients []ClientInfo = []ClientInfo{}
	for key, client := range clients {
		whoami := client.Store.ID.String()
		status := "disconnected"
		if client.IsConnected() {
			status = "connected"
		}
		clientInfo := ClientInfo{
			ID:      whoami,
			Phone:   whoami,
			Name:    key,
			Status:  status,
			Process: "getMessage",
			Busy:    true,
			Qrcode:  "",
		}
		connectedClients = append(connectedClients, clientInfo)
	}

	setResponse(w, 200, connectedClients)

}

func ScanDeviceHandler(w http.ResponseWriter, r *http.Request) {
	deviceStore := StoreContainer.NewDevice()
	client := GetClient(deviceStore)
	
	AddClient(client)

	// deviceID := GenerateRandomString("DEV", 7)
	// setClient_data(deviceID, client)
	
	qrCode, jid := connectClient(client)

	fmt.Println("data client")
	fmt.Println(data_client)

	var response []ClientInfo

	fmt.Println("Data client setelah ditambahkan:", jid)

	// Add the new client to the response
	if qrCode != "" {
		response = append(response, ClientInfo{
			ID:     "",
			Number: "",
			Busy:   false,
			QR:     qrCode,
			Status: "pairing",
			Name:   "",
		})
	}

	setResponse(w, 200, response)
}
