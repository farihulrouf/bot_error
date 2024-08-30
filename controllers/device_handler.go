package controllers

import (
	"fmt"
	"time"
	"reflect"
	"net/http"
	"wagobot.com/db"
	"wagobot.com/base"
	"wagobot.com/model"
)

func GetDevicesHandler(w http.ResponseWriter, r *http.Request) {
	// type ClientInfo struct {
	// 	ID      string `json:"id"`
	// 	Phone   string `json:"phone"`
	// 	Name    string `json:"name"`
	// 	Status  string `json:"status"`
	// 	Process string `json:"process"`
	// 	Busy    bool   `json:"busy"`
	// 	Qrcode  string `json:"qrcode"`
	// }

	mutex.Lock()
	defer mutex.Unlock()

	// tokenStr := r.Header.Get("Authorization")
	// tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	// claims, _ := base.ParseToken(tokenStr)
	// username, _ := claims["username"].(string)

	username := base.CurrentUser.Username
	user, _ := db.GetUserByUsername(username)

	// fmt.Println("User ID ", user.ID)

	var connectedClients []ClientInfo = []ClientInfo{}
	for _, client := range model.Clients {

		fmt.Println("id user", client)

		if reflect.ValueOf(client.Client.Store.ID).IsNil() {
			continue
		}
		
		if client.User != user.ID {
			continue
		}

		whoami := client.Client.Store.ID.String()
		phone := model.GetPhoneNumber(whoami)

		status := "disconnected"
		if client.Client.IsConnected() {
			status = "connected"
		}

		clientInfo := ClientInfo{
			ID:      whoami,
			// Phone:   phone,
			Number: phone,
			Name:    client.Client.Store.PushName,
			Status:  status,
			// Process: "getMessage",
			Busy:    true,
			QR:  "",
		}

		connectedClients = append(connectedClients, clientInfo)
	}

	base.SetResponse(w, http.StatusOK, connectedClients)
}

func ScanDeviceHandler(w http.ResponseWriter, r *http.Request) {

	// tokenStr := r.Header.Get("Authorization")
	// tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	// claims, _ := base.ParseToken(tokenStr)
	// username, _ := claims["username"].(string)

	username := base.CurrentUser.Username
	user, _ := db.GetUserByUsername(username)

	deviceStore := StoreContainer.NewDevice()
	client := GetClient(deviceStore)
	deviceID := GenerateRandomString("DEVICE", 5)

	currentTime := time.Now()
	nextTime := currentTime.Add(3 * time.Minute)
	nextUnixTime := nextTime.Unix()

	AddClient(user.ID, "", user.Url, deviceID, client, nextUnixTime)

	// fmt.Println(clients)

	qrCode, _ := connectClient(client)

	// fmt.Println("data client")
	// fmt.Println(data_client)

	var response ClientInfo

	// fmt.Println("Data client setelah ditambahkan:", jid)

	// Add the new client to the response
	if qrCode != "" {
		response = ClientInfo{
			ID:     "",
			Number: "",
			Busy:   false,
			QR:     qrCode,
			Status: "pairing",
			Name:   "",
		}
	}

	base.SetResponse(w, http.StatusOK, response)
}

func StatusDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var params model.PhoneRequest

	base.ValidateRequest(r, &params)
	fmt.Println(params)

	if params.Phone == "" {
		base.SetResponse(w, http.StatusBadRequest, "phone are required")
		return
	}

	phone := params.Phone

	if !base.IsMyNumber(phone) {
		base.SetResponse(w, http.StatusBadRequest, "Missing number")
		return
	}

	if _, exists := model.Clients[phone]; exists {
		client := model.Clients[phone].Client

		base.SetResponse(w, http.StatusOK, client.IsConnected())
	} else {
		base.SetResponse(w, http.StatusBadRequest, "Invalid account")
	}
}
