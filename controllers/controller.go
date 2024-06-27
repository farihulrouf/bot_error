package controllers

import (
	"context"

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
	"wagobot.com/db"
	"wagobot.com/model"
	"wagobot.com/response"
)

/*var client *whatsmeow.Client

func SetClient(c *whatsmeow.Client) {
	client = c
}

var (
	clients        = make(map[string]*whatsmeow.Client)
	mutex          = &sync.Mutex{}
	StoreContainer *sqlstore.Container
)

func SetStoreContainer(container *sqlstore.Container) {
	StoreContainer = container
}
*/

// maping client to map
var (
	clients        = make(map[string]*whatsmeow.Client)
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
}

func GetClients() map[string]*whatsmeow.Client {
	mutex.Lock()
	defer mutex.Unlock()
	return clients
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

func getClient(deviceStore *store.Device) *whatsmeow.Client {
	client := whatsmeow.NewClient(deviceStore, clientLog)
	//clients := client[generateRandomString(5)]
	//fmt.Println("data", clients)
	client.AddEventHandler(EventHandler)
	fmt.Println("data client", client)
	return client
}

/*
	func sendToAPI(sender string, message string) {
		mu.Lock()
		messages = append(messages, model.Message{Sender: sender, Message: message})
		mu.Unlock()
	}

*
//var silver = ""
*/
func EventHandler(evt interface{}) {
	fmt.Println("try to excution")
	switch v := evt.(type) {
	case *events.Message:
		if !v.Info.IsFromMe && v.Message.GetConversation() != "" ||
			v.Message.GetImageMessage().GetCaption() != "" ||
			v.Message.GetVideoMessage().GetCaption() != "" ||
			v.Message.GetDocumentMessage().GetCaption() != "" {
			id := v.Info.ID
			chat := v.Info.Sender.String()
			timestamp := v.Info.Timestamp
			text := v.Message.GetConversation()
			group := v.Info.IsGroup
			isfrome := v.Info.IsFromMe
			doc := v.Message.GetDocumentMessage()
			captionMessage := v.Message.GetImageMessage().GetCaption()
			videoMessage := v.Message.GetVideoMessage().GetCaption()
			docMessage := v.Message.GetDocumentMessage().GetCaption()
			docCaption := v.Message.GetDocumentMessage().GetTitle()
			name := v.Info.PushName
			to := v.Info.PushName

			//to : = v.Info.na
			thumbnail := v.Message.ImageMessage.GetJpegThumbnail()
			thumbnailvideo := v.Message.VideoMessage.GetJpegThumbnail()
			thumbnaildoc := v.Message.DocumentMessage.GetJpegThumbnail()
			url := v.Message.ImageMessage.GetUrl()
			mimeTipe := v.Message.ImageMessage.GetMimetype()
			//comment := v.Message.CommentMessage.GetMessage()
			//relyId := v.Message.GetCommentMessage().Message
			tipe := v.Info.Type
			isdocument := v.IsDocumentWithCaption
			//chatText := v.Info.Chat
			mediatype := v.Info.MediaType
			//smtext := v.Message.Conversation()
			fmt.Println("ID: %s, Chat: %s, Time: %d, Text: %s\n", to, mediatype, isdocument, chat, timestamp, text, group, isfrome, tipe)
			//fmt.Println("info repley", reply, coba)

			// Assuming replies are stored within a field named Replies
			fmt.Println("tipe messages", tipe, docCaption, isdocument, doc, mediatype, captionMessage, videoMessage, docMessage)
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
		//case *events.PairSuccess:
		//	fmt.Println("pari succeess", v.ID.User)
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

func ScanQrCode(client *whatsmeow.Client) {
	if client.Store.ID == nil {
		qrChannel, _ := client.GetQRChannel(context.Background())
		go func() {
			for evt := range qrChannel {
				switch evt.Event {
				case "code":
					fmt.Println("QR Code:", evt.Code)
				case "login":
					fmt.Println("Login successful")
				}
			}
		}()
		err := client.Connect()
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		<-qrChannel
	} else {
		err := client.Connect()
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
	}
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
		if textFilter != "" && !strings.Contains(msg.Text, textFilter) {
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
		data = append(data, messageData)
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

func CreateDevice(w http.ResponseWriter, r *http.Request) {
	deviceStore := StoreContainer.NewDevice()
	client := getClient(deviceStore)

	qrCode, jid := connectClient(client)
	fmt.Println("cek nilai client ", clients)

	var response []ClientInfo

	// Check if there are devices in the database
	dbx, err := db.OpenDatabase()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer dbx.Close()

	devices, err := model.GetDevicesFromDB(dbx)
	if err != nil {
		http.Error(w, "Failed to get devices from the database", http.StatusInternalServerError)
		return
	}

	// If there are devices, add them to the response
	if len(devices) > 0 {
		for _, d := range devices {
			response = append(response, ClientInfo{
				ID:     d.RegistrationID,
				Number: d.JID,
				Status: "connected",
				Name:   d.PushName,
			})
		}
	}

	// Add the new client to the response and clients map
	if qrCode != "" {
		response = append(response, ClientInfo{
			ID:     "",
			Number: "",
			QR:     qrCode,
			Status: "pairing",
			Name:   "",
		})
	}
	if jid != nil {
		clients[jid.User] = client
	}

	// Add existing connected clients to the response
	for _, c := range clients {
		if c.IsConnected() {
			response = append(response, ClientInfo{
				ID:     "93847384", // Example ID
				Number: c.Store.ID.User,
				Status: "connected",
				Name:   "dfarihul", // Example name
			})
		}
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
			"id":   msg.ID,
			"from": strings.TrimSuffix(msg.From, "@s.whatsapp.net"),
			"to":   msg.To,
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
		data = append(data, messageData)
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
			"id":   msg.ID,
			"chat": chat,
			"time": msg.Time,
			//"text": msg.Text,
			//"type": msg.Mediatipe,
		}
		if msg.Mediatipe != "text" {
			messageData["type"] = msg.Mediatipe
		}

		if msg.Tipe == "text" {
			messageData["text"] = msg.Text
		}

		// Tambahkan elemen ke slice
		data = append(data, messageData)
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

func RetrieveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	/*identifier := r.URL.Query().Get("identifier")
	if identifier == "" {
		http.Error(w, "Missing identifier", http.StatusBadRequest)
		return
	}

	messages, err := helpers.GetAllMessagesByPhoneNumberOrGroupID(client, identifier)
	if err != nil {
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	response := model.GetMessagesResponse{Data: messages}
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

func TriggerEventHandler(w http.ResponseWriter, r *http.Request) {

	for _, client := range clients {
		client.AddEventHandler(EventHandler)
	}
	for key := range clients {
		fmt.Println("Checking key:", key)
		EventHandler(clients[key])
	}
}

func ListDevices(w http.ResponseWriter, r *http.Request) {
	db, err := db.OpenDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	devices, err := model.GetDevicesFromDB(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve devices: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
