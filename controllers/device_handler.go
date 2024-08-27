package controllers

import (
	// "encoding/json"
	"fmt"
	// "reflect"
	"strings"
	"net/http"
	"wagobot.com/auth"
	"wagobot.com/db"
	"wagobot.com/model"
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

	tokenStr := r.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, _ := auth.ParseToken(tokenStr)
	username, _ := claims["username"].(string)

	user, _ := db.GetUserByUsername(username)

	fmt.Println("User ID ", user.ID)

	var connectedClients []ClientInfo = []ClientInfo{}
	for _, client := range clients {

		// fmt.Println("id user", client)
		
		if client.User != user.ID {
			continue
		}

		fmt.Println(client.Client)

		whoami := client.Client.Store.ID.String()
		phone := model.GetPhoneNumber(whoami)
		status := "disconnected"
		if client.Client.IsConnected() {
			status = "connected"
		}
		clientInfo := ClientInfo{
			ID:      whoami,
			Phone:   phone,
			Name:    client.Client.Store.PushName,
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

	tokenStr := r.Header.Get("Authorization")
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	claims, _ := auth.ParseToken(tokenStr)
	username, _ := claims["username"].(string)

	user, _ := db.GetUserByUsername(username)

	fmt.Println("User ID ", user.ID)

	deviceStore := StoreContainer.NewDevice()
	client := GetClient(deviceStore)
	
	deviceID := GenerateRandomString("DEVICE", 5)

	AddClient(user.ID, deviceID, client)

	fmt.Println(clients)

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
