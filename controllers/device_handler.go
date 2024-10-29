package controllers

import (
	"fmt"
	"net/http"
	"time"

	"wagobot.com/base"
	"wagobot.com/db"
	"wagobot.com/model"
)

func ScanDeviceHandler(w http.ResponseWriter, r *http.Request) {
	// var params model.PhoneRefParams

	params := r.URL.Query()
	phone := params.Get("phone")
	ref := params.Get("ref")

	fmt.Println(params)

	if phone == "" || ref == "" {
		base.SetResponse(w, http.StatusBadRequest, "phone and ref required")
		return
	}

	username := base.CurrentUser.Username
	user, _ := db.GetUserByUsername(username)

	deviceStore := StoreContainer.NewDevice()
	client := GetClient(deviceStore)
	// deviceID := GenerateRandomString("DEVICE", 5)
	deviceID := "DEVICE" + "-" + phone + "-" + ref

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

func RemoveDeviceHandler(w http.ResponseWriter, r *http.Request) {
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
		// client := model.Clients[phone].Client

		model.Clients[phone].Client.Logout()
		delete(model.Clients, phone)

		base.SetResponse(w, http.StatusOK, true)
	} else {
		base.SetResponse(w, http.StatusBadRequest, "Invalid account")
	}
}
