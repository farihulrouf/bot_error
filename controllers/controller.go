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
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"wagobot.com/errors"
	"wagobot.com/helpers"
	"wagobot.com/model"
	"wagobot.com/response"
)

var client *whatsmeow.Client

func SetClient(c *whatsmeow.Client) {
	client = c
}

var (
	messages   []response.Message
	mu         sync.Mutex
	webhookURL string
)

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

func ScanQRHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the client is nil
	if client == nil {
		http.Error(w, "Client is nil", http.StatusInternalServerError)
		return
	}

	// No ID stored, new login
	qrChan, _ := client.GetQRChannel(context.Background())
	err := client.Connect()
	if err != nil {
		http.Error(w, "Failed to connect: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Loop through QR channel events
	for evt := range qrChan {
		if evt.Event == "code" {
			// Respond with the QR code data
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"qr_code": evt.Code})
			return
		}
	}
}

func RetrieveMessagesHandler(w http.ResponseWriter, r *http.Request) {
	identifier := r.URL.Query().Get("identifier")
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
}

/*
func GetDevicesHandler(w http.ResponseWriter, r *http.Request) {
	// Get self JID from the device store
	deviceStore := client.Store.ID
	if deviceStore == nil {
		http.Error(w, "Client not logged in", http.StatusInternalServerError)
		return
	}

	// Convert the deviceStore ID to a proper JID
	selfJID := types.NewJID(deviceStore.User, types.DefaultUserServer)

	// Get user devices for the logged-in JID
	deviceJIDs, err := client.GetUserDevices([]types.JID{selfJID})
	if err != nil {
		log.Printf("Error getting user devices: %v", err)
		http.Error(w, "Failed to get user devices", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	responseData := make([]map[string]interface{}, 0)
	for _, deviceJID := range deviceJIDs {
		// Fetch user info for each device
		userInfoMap, err := client.GetUserInfo([]types.JID{deviceJID})
		if err != nil {
			log.Printf("Error getting user info: %v", err)
			// Include the device in the response with limited information
			deviceData := map[string]interface{}{
				"id":      deviceJID.String(),
				"phone":   deviceJID.User,
				"status":  "unknown",
				"process": "string", // Replace with actual process if available
				"busy":    false,    // Replace with actual busy status if available
				"qrcode":  "",       // Replace with actual QR code if available
			}
			responseData = append(responseData, deviceData)
			continue // Continue to the next device
		}

		userInfo, exists := userInfoMap[deviceJID]
		if !exists {
			http.Error(w, "User info not found for device", http.StatusNotFound)
			return
		}

		fmt.Println("check userInfo", userInfo)

		deviceData := map[string]interface{}{
			"id":    deviceJID.String(),
			"phone": deviceJID.User,
			//"name":    userInfo.Long, // Use Long name instead of Short
			"status":  userInfo.Status,
			"process": "string", // Replace with actual process if available
			"busy":    false,    // Replace with actual busy status if available
			"qrcode":  "",       // Replace with actual QR code if available
		}

		fmt.Println("cek data", deviceData)
		responseData = append(responseData, deviceData)
	}

	response := map[string]interface{}{
		"data": responseData,
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
*/
/*
func SetWebhook(w http.ResponseWriter, r *http.Request) {
	txtid := r.Context().Value("userinfo").(auth.Values).Get("Id")
	token := r.Context().Value("userinfo").(auth.Values).Get("Token")
	userid, _ := strconv.Atoi(txtid)

	decoder := json.NewDecoder(r.Body)
	var t model.WebhookStruct
	err := decoder.Decode(&t)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("Could not set webhook: %v", err)))
		return
	}
	var webhook = t.WebhookURL

	_, err = s.db.Exec("UPDATE users SET webhook=? WHERE id=?", webhook, userid)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, errors.New(fmt.Sprintf("%s", err)))
		return
	}

	v := helpers.UpdateUserInfo(r.Context().Value("userinfo"), "Webhook", webhook)
	userinfocache.Set(token, v, cache.NoExpiration)

	response := map[string]interface{}{"webhook": webhook}
	responseJson, err := json.Marshal(response)
	if err != nil {
		Respond(w, r, http.StatusInternalServerError, err)
	} else {
		Respond(w, r, http.StatusOK, string(responseJson))
	}
	return
}

func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
*/
