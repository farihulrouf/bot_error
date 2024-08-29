package controllers

import (
	"context"
	"io/ioutil"
	"math/rand"
	"time"

	// "encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	// "go.mau.fi/whatsmeow/types/events"
	
	"wagobot.com/model"
	"wagobot.com/response"
	// "wagobot.com/db"
	"wagobot.com/base"

	waLog "go.mau.fi/whatsmeow/util/log"
)

// maping client to map
const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type GroupCollection struct {
	Groups []types.GroupInfo
}

var (
	// clients        = make(map[string]*whatsmeow.Client)
	// clients        = make(map[string]CustomClient)
	data_client    = make(map[string]*whatsmeow.Client)
	mutex          = &sync.Mutex{}
	StoreContainer *sqlstore.Container
	clientLog      waLog.Logger
)

// set store
func SetStoreContainer(container *sqlstore.Container) {
	StoreContainer = container
}

var (
	messages   []response.Message
	mu         sync.Mutex
	webhookURL string
)

type ClientInfo struct {
	ID     string `json:"id"`
	Number string `json:"number,omitempty"`
	QR     string `json:"qr,omitempty"`
	Status string `json:"status"`
	Name   string `json:"name"`
	Busy   bool   `json:"busy,omitempty"`
}

func CleanupClients() {
	fmt.Println("Cleanup", model.Clients)
	currentTime := time.Now()
	currentUnixTime := currentTime.Unix()
	for key, client := range model.Clients {
		if !strings.HasPrefix(key, "DEV") {
			continue
		}
		// fmt.Println("Expired time ", client.ExpiredTime, currentUnixTime)
		if client.ExpiredTime > 0 && client.ExpiredTime < currentUnixTime {
			// fmt.Println("ini expired ", key)
			delete(model.Clients, key)
		}
	}
}

func setClient_data(key string, client *whatsmeow.Client) {
	// Clear existing data
	for k := range data_client {
		delete(data_client, k)
	}
	// Set new client
	data_client[key] = client
}

func connectClient(client *whatsmeow.Client) (string, *types.JID) {
	var err error
	qrChan := make(chan string)

	// Disconnect client if it's already connected
	if client.IsConnected() {
		client.Disconnect()
	}

	// Generate new QR code for new login session
	qrChannel, _ := client.GetQRChannel(context.Background())
	go func() {
		for evt := range qrChannel {
			switch evt.Event {
			case "code":
				qrChan <- evt.Code
			case "login":
				close(qrChan)
			}
		}
	}()

	err = client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	qrCode := <-qrChan
	return qrCode, client.Store.ID
}

func GetClient(deviceStore *store.Device) *whatsmeow.Client {
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	// handler := &CustomEventHandler{client: client}
	// var _ whatsmeow.EventHandler = handler
	// client.AddEventHandler(handler)
	client.AddEventHandler(EventHandler)
	return client
}

func GetSearchMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters
	r.ParseForm()
	textFilter := r.Form.Get("chat")

	mu.Lock()
	defer mu.Unlock()
	//fmt.Println("check data message five", v.Info.IsGroup, v.Info.IsFromMe, v.Info.Category, v.Info.MessageSource, v.Info.Type, v.Info.Chat.Device, v.Info.Timestamp)

	w.Header().Set("Content-Type", "application/json")
	var data []map[string]interface{}
	//data := make(map[string]map[string]interface{})
	for _, msg := range messages {
		if textFilter != "" && !strings.Contains(msg.Text, textFilter) &&
			!strings.Contains(msg.Chat, textFilter) &&
			!strings.Contains(msg.ID, textFilter) &&
			!strings.Contains(msg.Name, textFilter) &&
			!strings.Contains(msg.ID, textFilter) &&
			!strings.Contains(msg.Caption, textFilter) {
			continue // Skip messages that don't contain the text filter
		}

		//timeStr := fmt.Sprintf("%d", msg.Time)
		// Remove @s.whatsapp.net suffix from msg.Chat
		chat := strings.TrimSuffix(msg.Chat, "@s.whatsapp.net")
		//tiipe chat group or user
		chatType := "user"
		thumb := msg.Thumbnail
		if msg.Group {
			chatType = "Group"
		}
		if msg.Tipe != "text" {
			msg.Text = msg.Caption
		}

		if msg.Mediatipe == "image" {
			msg.Text = msg.Caption
			thumb = msg.Thumbnail

		}
		if msg.Mediatipe == "video" {
			msg.Text = msg.VideoMessage
			thumb = msg.Thumbnailvideo
		}
		if msg.Mediatipe == "document" {
			msg.Text = msg.DocMessage
			thumb = msg.Thumbnaildoc
		}
		if msg.Mediatipe == "" {
			msg.Mediatipe = "text"
		}

		messageData := map[string]interface{}{
			"id":        msg.ID,
			"time":      msg.Time,
			"fromMe":    true, //!v.Info.IsFromMe && v.Message.GetConversation() !=
			"type":      msg.Mediatipe,
			"status":    "delivered",
			"chatType":  chatType,
			"replyId":   "1609773514305",
			"chat":      chat,
			"to":        chat,
			"name":      msg.Name,
			"from":      chat,
			"text":      msg.Text,
			"caption":   msg.Caption,
			"url":       msg.Url,
			"mimetype":  msg.MimeTipe,
			"thumbnail": thumb,
		}

		exists := false
		for _, existingMessage := range data {
			if existingMessage["id"] == msg.ID {
				exists = true
				break
			}
		}

		// Jika msg.ID belum ada, tambahkan messageData ke data
		if !exists {
			data = append(data, messageData)
		}

		//data = append(data, messageData)
		//fmt.Println("chek data", msg)
		/* example respond in maxchat.id
		{
			"data": [
				{
				"id": "1609773514305",
				"time": 1686380234054,
				"fromMe": true,
				"type": "text",
				"status": "delivered",
				"chatType": "user",
				"replyId": "1609773514305",
				"chat": "6281234567890",
				"to": "6281234567890",
				"name": "string",
				"from": "6281234567890",
				"text": "Test from MaxChat",
				"caption": "Caption test",
				"url": "https://www.fnordware.com/superpng/pnggrad16rgb.png",
				"mimetype": "string",
				"thumbnail": "string"
				}
			]
			}
		*/
		//data[timeStr] = messageData
	}

	response := map[string]interface{}{
		"data": data,
	}

	// Encode response to JSON and send it
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
}

func AddClient(UserID int, DevID string, client *whatsmeow.Client, expired int64) {
	mutex.Lock()
	defer mutex.Unlock()

	if client == nil {
		log.Printf("Failed to add client: client is nil")
		return
	}

	// handler := &CustomEventHandler{client: client}
	// client.AddEventHandler(handler)
	client.AddEventHandler(EventHandler)

	// devId := GenerateRandomString("DEVICE", 5)
	// if _, ok := clients[devId]; !ok {
		model.Clients[DevID] = model.CustomClient{
			User: UserID,
			ExpiredTime: expired,
			Client: client,
		}
	// }

	err := model.Clients[DevID].Client.Connect()
	if err != nil {
		log.Fatalf("Gagal menghubungkan klien: %v", err)
	}

	// clients[id] = client
	log.Printf("Client added successfully: %s\n", DevID)
	fmt.Println(model.Clients)
}

func CreateDevice(w http.ResponseWriter, r *http.Request) {
	deviceStore := StoreContainer.NewDevice()
	client := GetClient(deviceStore)
	deviceID := GenerateRandomString("Device", 3)
	//data_client[deviceID] = client
	setClient_data(deviceID, client)
	qrCode, jid := connectClient(client)

	var response []ClientInfo

	fmt.Println("Data client setelah ditambahkan:", jid)

	// Iterasi melalui peta `clients` untuk membuat respons
	for key, client := range model.Clients {
		//fmt.Printf(key)
		response = append(response, ClientInfo{
			ID:     key,
			Number: client.Client.Store.ID.String(),
			Busy:   true,
			QR:     "",
			Status: "connected",
			Name:   client.Client.Store.PushName,
		})
	}

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

	w.Header().Set("Content-Type", "application/json")
	if len(response) > 0 {
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Failed to connect the client", http.StatusInternalServerError)
	}
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	// Ubah dari map ke slice
	var data []map[string]interface{}

	for _, msg := range messages {
		/*if msg.Mediatipe == "image" {
			msg.Text = msg.Caption
		}
		if msg.Mediatipe == "video" {
			msg.Text = msg.VideoMessage
		}
		if msg.Mediatipe == "document" {
			msg.Text = msg.DocMessage
		}
		if msg.Mediatipe == "" {
			msg.Mediatipe = "text"
		}
		*/
		/*
			if msg.DocMessage != "" {
				msg.Text = msg.DocMessage
			}
		*/

		//chat := strings.TrimSuffix(msg.Chat, "@s.whatsapp.net"Ev)

		if msg.Tipe == "text" {
			msg.Mediatipe = "text"
		}
		messageData := map[string]interface{}{
			"id":     msg.ID,
			"from":   strings.TrimSuffix(msg.From, "@s.whatsapp.net"),
			"to":     strings.TrimSuffix(msg.From, "@s.whatsapp.net"),
			"status": "delivered",
			//"chat": chat,
			"time": msg.Time,
			"type": msg.Mediatipe,
			//"text": msg.Text,
		}
		/*if msg.Tipe != "text" {
			messageData["type"] = msg.Mediatipe
		}
		*/
		if msg.Tipe == "text" {
			messageData["text"] = msg.Text
		}

		// Tambahkan elemen ke slice
		exists := false
		for _, existingMessage := range data {
			if existingMessage["id"] == msg.ID {
				exists = true
				break
			}
		}

		// Jika msg.ID belum ada, tambahkan messageData ke data
		if !exists {
			data = append(data, messageData)
		}
	}

	response := map[string]interface{}{
		"data": data,
	}

	// Encode response ke JSON dan kirim
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
}

func GetMessagesByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil id dari URL
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	// Ubah dari map ke slice
	var data []map[string]interface{}

	for _, msg := range messages {
		// Remove @s.whatsapp.net suffix from msg.Chat
		chat := strings.TrimSuffix(msg.Chat, "@s.whatsapp.net")

		// Filter pesan berdasarkan chat
		if chat != id {
			continue
		}

		/*if msg.Mediatipe == "image" {
			msg.Text = msg.Caption
		}
		if msg.Mediatipe == "video" {
			msg.Text = msg.VideoMessage
		}
		if msg.Mediatipe == "" {
			msg.Mediatipe = "text"
		}
		if msg.Mediatipe == "document" {
			msg.Text = msg.DocMessage
		}
		*/

		messageData := map[string]interface{}{
			"id":     msg.ID,
			"chat":   chat,
			"time":   msg.Time,
			"status": "delivered",
			//"text": msg.Text,
			//"type": msg.Mediatipe,
		}
		if msg.Mediatipe != "text" {
			messageData["type"] = msg.Mediatipe
		}

		if msg.Tipe == "text" {
			messageData["text"] = msg.Text
		}

		exists := false
		for _, existingMessage := range data {
			if existingMessage["id"] == msg.ID {
				exists = true
				break
			}
		}

		if !exists {
			data = append(data, messageData)
		}

	}

	response := map[string]interface{}{
		"data": data,
	}

	// Encode response ke JSON dan kirim
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
}

func initialClient() {
	for key, value := range data_client {
		model.Clients[key] = model.CustomClient {
			User: 0,
			Client: value,
		}
	}
}

func GenerateRandomString(prefix string, length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("%s-%s", prefix, string(b))
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received raw payload: %s\n", string(body)) // Logging payload mentah

	var payload model.WebhookPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Failed to parse webhook payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed webhook payload: %+v\n", payload)

	w.WriteHeader(http.StatusOK)
}

func RemoveClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	phone := vars["phone"]

	// Lock untuk mengamankan akses ke map clients (jika diperlukan)
	mutex.Lock()
	defer mutex.Unlock()

	// Cek apakah kunci ada di dalam map
	if _, exists := model.Clients[phone]; exists {
		// Hapus kunci dari map
		model.Clients[phone].Client.Logout()
		delete(model.Clients, phone)
		delete(data_client, phone)

		// response := response.ResponseLogout{Status: "success", Message: "Data berhasil dihapus"}
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(response)

		base.SetResponse(w, http.StatusOK, "Data berhasil dihapus")
	} else {
		// response := response.ResponseLogout{Status: "fail", Message: "Kunci tidak ditemukan"}
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusNotFound)
		// json.NewEncoder(w).Encode(response)

		base.SetResponse(w, http.StatusNotFound, "Kunci tidak ditemukan")
	}

}
