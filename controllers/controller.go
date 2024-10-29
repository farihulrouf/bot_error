package controllers

import (
	"context"
	"math/rand"
	"time"

	// "encoding/base64"

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
		if strings.HasPrefix(key, "DEV") {
			// fmt.Println("Expired time ", client.ExpiredTime, currentUnixTime)
			if client.ExpiredTime > 0 && client.ExpiredTime < currentUnixTime {
				// fmt.Println("ini expired ", key)
				delete(model.Clients, key)
			}
		} else {
			// untuk nomor dengan id user = 0
			if client.User == 0 {
				// logout lalu hapus
				if client.Client.IsLoggedIn() {
					client.Client.Logout()
				}
				delete(model.Clients, key)
			}
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
	// client.AddEventHandler(EventHandler)
	client.AddEventHandler(func(evt interface{}) {
		EventHandler(evt, model.CustomClient{
			Client: client,
		})
	})
	return client
}

func AddClient(UserID int, phone string, webhook string, DevID string, client *whatsmeow.Client, expired int64) {
	mutex.Lock()
	defer mutex.Unlock()

	if client == nil {
		log.Printf("Failed to add client: client is nil")
		return
	}

	customClient := model.CustomClient{
		User:        UserID,
		Phone:       phone,
		ExpiredTime: expired,
		Webhook:     webhook,
		Client:      client,
	}

	// handler := &CustomEventHandler{client: client}
	// client.AddEventHandler(handler)
	// client.AddEventHandler(EventHandler)
	client.AddEventHandler(func(evt interface{}) {
		EventHandler(evt, customClient)
	})

	model.Clients[DevID] = customClient

	// devId := GenerateRandomString("DEVICE", 5)
	// if _, ok := clients[devId]; !ok {
	// model.Clients[DevID] = model.CustomClient{
	// 	User: UserID,
	// 	ExpiredTime: expired,
	// 	Webhook: webhook,
	// 	Client: client,
	// }
	// }

	err := model.Clients[DevID].Client.Connect()
	if err != nil {
		log.Fatalf("Gagal menghubungkan klien: %v", err)
	}

	// clients[id] = client
	log.Printf("Client added successfully: %s\n", DevID)
	fmt.Println(model.Clients)
}

func initialClient() {
	for key, value := range data_client {
		model.Clients[key] = model.CustomClient{
			User:   0,
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
