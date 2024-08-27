package controllers

import (
	"context"
	"io/ioutil"
	"math/rand"
	"time"

	"encoding/base64"
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
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"wagobot.com/model"
	"wagobot.com/response"
)

// maping client to map
const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	clients        = make(map[string]*whatsmeow.Client)
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
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(EventHandler)
	return client
}

func EventHandler(evt interface{}) {

	switch v := evt.(type) {
	case *events.Message:
		for key := range clients {
			//fmt.Println("whoami:", key)
			fmt.Println("end client", clients[key])
		}

		if !v.Info.IsFromMe || v.Message.GetConversation() != "" ||
			v.Message.GetImageMessage().GetCaption() != "" ||
			v.Message.GetVideoMessage().GetCaption() != "" ||
			v.Message.GetDocumentMessage().GetCaption() != "" {
			id := v.Info.ID
			chat := v.Info.Sender.String()
			timestamp := v.Info.Timestamp
			text := v.Message.GetConversation()
			group := v.Info.IsGroup
			isfrome := v.Info.IsFromMe
			//doc := v.Message.GetDocumentMessage()
			captionMessage := v.Message.GetImageMessage().GetCaption()
			videoMessage := v.Message.GetVideoMessage().GetCaption()
			docMessage := v.Message.GetDocumentMessage().GetCaption()
			//docCaption := v.Message.GetDocumentMessage().GetTitle()
			name := v.Info.PushName
			to := v.Info.PushName

			//to : = v.Info.na
			//thumbnail := v.Message.ImageMessage.GetJpegThumbnail()
			thumbnail := v.Message.ImageMessage.GetJPEGThumbnail()
			thumbnailvideo := v.Message.VideoMessage.GetJPEGThumbnail()
			thumbnaildoc := v.Message.DocumentMessage.GetJPEGThumbnail()
			url := v.Message.ImageMessage.GetURL()
			mimeTipe := v.Message.ImageMessage.GetMimetype()

			tipe := v.Info.Type
			isdocument := v.IsDocumentWithCaption
			//chatText := v.Info.Chat
			mediatype := v.Info.MediaType
			//smtext := v.Message.Conversation()
			//fmt.Println("ID: %s, Chat: %s, Time: %d, Text: %s\n", to, mediatype, isdocument, chat, timestamp, text, group, isfrome, tipe)
			//fmt.Println("info repley", reply, coba)

			// Assuming replies are stored within a field named Replies
			//fmt.Println("tipe messages", tipe, docCaption, isdocument, doc, mediatype, captionMessage, videoMessage, docMessage)
			mu.Lock()
			defer mu.Unlock() // Ensure mutex is always unlocked when the function returns
			messages = append(messages, response.Message{
				ID:             id,
				Chat:           chat,
				Time:           timestamp.Unix(),
				Text:           text,
				Group:          group,
				Mediatipe:      mediatype,
				IsDocument:     isdocument,
				Tipe:           tipe,
				IsFromMe:       isfrome,
				Caption:        captionMessage,
				VideoMessage:   videoMessage,
				DocMessage:     docMessage,
				Name:           name,
				From:           chat,
				To:             to,
				Url:            url,
				Thumbnail:      base64.StdEncoding.EncodeToString(thumbnail),
				MimeTipe:       mimeTipe,
				Thumbnaildoc:   base64.StdEncoding.EncodeToString(thumbnaildoc),
				Thumbnailvideo: base64.StdEncoding.EncodeToString(thumbnailvideo),
				//MimeType:     *mimesType,
				//CommentMessage: comment,
				//Replies: reply,
				// Add replies to the message if available
				// Replies: v.Message.Replies,
			})
		}
		/*
			payload := response.Message{
				ID:             v.Info.ID,
				Chat:           v.Info.Sender.String(),
				Time:           v.Info.Timestamp.Unix(),
				Text:           v.Message.GetConversation(),
				Group:          v.Info.IsGroup,
				IsFromMe:       v.Info.IsFromMe,
				Caption:        v.Message.GetImageMessage().GetCaption(),
				VideoMessage:   v.Message.GetVideoMessage().GetCaption(),
				DocMessage:     v.Message.GetDocumentMessage().GetCaption(),
				MimeTipe:       v.Message.GetImageMessage().GetMimetype(),
				Name:           v.Info.PushName,
				To:             v.Info.PushName,
				Url:            v.Message.GetImageMessage().GetUrl(),
				Thumbnail:      base64.StdEncoding.EncodeToString(v.Message.GetImageMessage().GetJpegThumbnail()),
				Thumbnailvideo: base64.StdEncoding.EncodeToString(v.Message.GetVideoMessage().GetJpegThumbnail()),
				Thumbnaildoc:   base64.StdEncoding.EncodeToString(v.Message.GetDocumentMessage().GetJpegThumbnail()),
				Tipe:           v.Info.Type,
				IsDocument:     v.IsDocumentWithCaption,
				Mediatipe:      v.Info.MediaType,
			}
			webhookURL := "https://webhook.site/aa9bbb63-611c-4d7a-97cd-f4eb6d4b775d"
			err := sendPayloadToWebhook(payload, webhookURL)
			if err != nil {
				fmt.Printf("Failed to send payload to webhook: %v\n", err)
			}
		*/
	case *events.PairSuccess:
		fmt.Println("pari succeess", v.ID.User)
		initialClient()
	case *events.HistorySync:
		fmt.Println("Received a history sync")
		/*for _, conv := range v.Data.GetConversations() {
			for _, historymsg := range conv.GetMessages() {
				chatJID, _ := types.ParseJID(conv.GetId())
				evt, err := client.ParseWebMessage(chatJID, historymsg.GetMessage())
				if err != nil {
					log.Println(err)
				}
				eventHandler(evt)
			}
		}*/
	case *events.LoggedOut:
		//initialClient()
	case *events.Receipt:
		fmt.Println("terima")
		if v.Type == events.ReceiptTypeRead || v.Type == events.ReceiptTypeReadSelf {
			fmt.Printf("%v was read by %s at %s\n", v.MessageIDs, v.SourceString(), v.Timestamp)
			// Membuat payload untuk webhook
			/*webhookPayload := model.ReadReceipt{
				MessageID: v.MessageIDs[0],
				ReadBy:    v.SourceString(),
				Time:      v.Timestamp.UnixMilli(),
			}*/
			// Mengirimkan payload ke webhook
			//webhookURL := "http://localhost:8080/webhook"
			/*err := sendPayloadToWebhook(string(v.Type), webhookURL)
			if err != nil {
				fmt.Printf("Failed to send read receipt to webhook: %v\n", err)
			}*/
		}

		/*case *events.Receipt:
		if v.Type == types.ReceiptTypeRead || v.Type == types.ReceiptTypeReadSelf {
			fmt.Println("%v was read by %s at %s", v.MessageIDs, v.SourceString(), v.Timestamp)
		} else if v.Type == types.ReceiptTypeDelivered {
			fmt.Println("%s was delivered to %s at %s", v.MessageIDs[0], v.SourceString(), v.Timestamp)
		}
		*/

	}
}

type GroupCollection struct {
	Groups []types.GroupInfo
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

func AddClient(id string, client *whatsmeow.Client) {
	mutex.Lock()
	defer mutex.Unlock()

	if client == nil {
		log.Printf("Failed to add client: client is nil for id %s\n", id)
		return
	}

	clients[id] = client
	log.Printf("Client added successfully: %s\n", id)
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
	for key, client := range clients {
		//fmt.Printf(key)
		response = append(response, ClientInfo{
			ID:     key,
			Number: client.Store.ID.String(),
			Busy:   true,
			QR:     "",
			Status: "connected",
			Name:   client.Store.PushName,
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
		clients[key] = value
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
	if _, exists := clients[phone]; exists {
		// Hapus kunci dari map
		clients[phone].Logout()
		delete(clients, phone)
		delete(data_client, phone)

		response := response.ResponseLogout{Status: "success", Message: "Data berhasil dihapus"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := response.ResponseLogout{Status: "fail", Message: "Kunci tidak ditemukan"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
	}

}
