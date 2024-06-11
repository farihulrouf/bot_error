package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
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
*/
func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if !v.Info.IsFromMe && v.Message.GetConversation() != "" {
			id := v.Info.ID
			chat := v.Info.Sender.String()
			timestamp := time.Now().UnixNano() / int64(time.Millisecond)
			text := v.Message.GetConversation()
			//reply := v.Message.ReactionMessage
			//coba := v.Message.DeviceSentMessage
			fmt.Printf("ID: %s, Chat: %s, Time: %d, Text: %s\n", id, chat, timestamp, text)
			//fmt.Println("info repley", reply, coba)

			// Assuming replies are stored within a field named Replies

			mu.Lock()
			defer mu.Unlock() // Ensure mutex is always unlocked when the function returns
			messages = append(messages, response.Message{
				ID:   id,
				Chat: chat,
				Time: timestamp,
				Text: text,

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

func GetGroupsHandler(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	groups, err := client.GetJoinedGroups()
	if err != nil {
		http.Error(w, "Failed to fetch joined groups", http.StatusInternalServerError)
		return
	}
	var filteredGroups []response.GroupResponse
	for _, group := range groups {
		var isMember bool
		var members []string
		var admins []string

		for _, member := range group.Participants {
			if member.JID.User == phone {
				isMember = true
			}
			members = append(members, member.JID.User)

			// Check if the member is an admin
			if member.IsAdmin {
				admins = append(admins, member.JID.User)
			}
		}

		if isMember {
			groupID := strings.TrimSuffix(group.JID.String(), "@g.us")
			groupResponse := response.GroupResponse{
				ID:          groupID,
				Type:        "group",
				Description: group.Name,
				Members:     members,
				Admins:      admins, //time.Now().UnixMilli()
				Time:        group.GroupCreated.UnixMilli(),
				Pinned:      false,
				UnreadCount: 30,
			}
			filteredGroups = append(filteredGroups, groupResponse)
		}
	}

	response := map[string]interface{}{
		"data": filteredGroups,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure all required fields are present
	if req.Code == "" || req.Phone == "" {
		http.Error(w, "code and phone are required fields", http.StatusBadRequest)
		return
	}

	if !helpers.IsLoggedInByNumber(client, req.Phone) {
		sendErrorResponse(w, http.StatusBadRequest, "not ready or not available. Please pairing the device", req.Phone)
		return
	}

	// Attempt to join the group with the provided invite link (code)
	groupJID, err := client.JoinGroupWithLink(req.Code)
	if err != nil {
		log.Printf("Error joining group: %v", err)
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	// Log success and respond with a success message
	log.Printf("Group joined successfully: %v", groupJID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group joined successfully"})
}

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LeaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	/*
		if !helpers.IsLoggedInByNumber(client, requestData.From) {
			http.Error(w, "Sender number is not connected to WhatsApp", http.StatusBadRequest)
			return
		}
	*/

	groupJID, err := types.ParseJID(req.GroupID + "@g.us")
	if err != nil {
		log.Printf("Error parsing Group JID: %v", err)
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	participantJID, err := types.ParseJID(req.Phone + "@s.whatsapp.net")
	if err != nil {
		log.Printf("Error parsing Participant JID: %v", err)
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	//  Dengan asumsi metode untuk memperbarui peserta grup adalah UpdateGroupParticipants
	response, err := client.UpdateGroupParticipants(groupJID, []types.JID{participantJID}, "remove")
	if err != nil {
		log.Printf("Error leaving group: %v", err)
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	log.Printf("Participant %s left group %s successfully, response: %v", req.Phone, req.GroupID, response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"})
}

func SendMessageGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.SendMessageDataRequest

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the request data
	if req.Type == "" || req.Text == "" {
		http.Error(w, "Missing required fields: 'type' and 'text'", http.StatusBadRequest)
		return
	}

	// Convert to JID
	jid, err := helpers.ConvertToJID(req.To)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid recipient: %v", err), http.StatusBadRequest)
		return
	}

	// Send the message
	if err := helpers.SendMessage(client, jid, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message sent to: %s", req.To)
}

// SendMessageHandler handles sending messages.

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the message data
	//var isAdmin bool
	//adminGroupJIDs := make([]string, 0)

	var requestData model.SendMessageDataRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if any required field is missing
	if requestData.To == "" || requestData.Type == "" || requestData.Text == "" || requestData.From == "" {
		http.Error(w, "Missing required fields: 'to', 'type', 'text', or 'from'", http.StatusBadRequest)
		return
	}

	// Validate phone numbers
	if !helpers.IsValidPhoneNumber(requestData.To) {
		http.Error(w, "Invalid phone number to", http.StatusBadRequest)
		return
	}
	if !helpers.IsValidPhoneNumber(requestData.From) {
		http.Error(w, "Invalid phone number sender", http.StatusBadRequest)
		return
	}
	//Check numerphone is login or not
	if !helpers.IsLoggedInByNumber(client, requestData.From) {
		sendErrorResponse(w, http.StatusBadRequest, "not ready or not available. Please pairing the device", requestData.From)
		return
	}
	// Simpan daftar JID grup yang merupakan admin

	// Periksa setiap grup untuk memeriksa apakah pengguna adalah admin
	/* tidak perlu mengirim pesan ke group hanya private / only private
	groups, err := client.GetJoinedGroups()
	if err != nil {
		http.Error(w, "Failed to fetch joined groups", http.StatusInternalServerError)
		return
	}

	///looping for message to group
	for _, group := range groups {
		for _, participant := range group.Participants {
			if participant.JID.User == requestData.To && participant.IsAdmin {
				// Jika nomor tersebut adalah admin, tambahkan JID grup ke dalam slice
				adminGroupJIDs = append(adminGroupJIDs, group.JID.String())

				//fmt.Println("check admin", adminGroupJIDs)
				// Keluar dari loop inner karena sudah ditentukan bahwa nomor tersebut adalah admin dalam grup ini
				break
			}
		}
	}

	fmt.Println("nomer jid", adminGroupJIDs)
	for _, groupJID := range adminGroupJIDs {
		// Convert the string JID to types.JID
		//jid := types.JID(groupJID)
		parts := strings.Split(groupJID, "@")

		// Extract user and server parts
		user := parts[0]
		server := parts[1]

		// Convert the user and server parts to types.JID
		jid := types.NewJID(user, server)

		// Call the SendMessageToGroup function
		err := helpers.SendMessageToGroup(client, jid, requestData.Text)
		if err != nil {
			fmt.Printf("Error sending message to group %s: %v\n", groupJID, err)
		} else {
			fmt.Printf("Message sent to group %s successfully\n", groupJID)
		}
	}
	*/
	if requestData.Type == "text" {
		err = helpers.SendMessageToPhoneNumber(client, requestData.To, requestData.Text)
		if err != nil {
			// Tangani kesalahan jika gagal mengirim pesan
			fmt.Printf("Error sending message to number", err)
		}
	} else {
		http.Error(w, "Invalid message type", http.StatusBadRequest)
		return
	}
	//fmt.Println("check admin xsilver", adminGroupJIDs)

	response := model.SendMessageResponse{
		ID:     uuid.New().String(),
		From:   requestData.From,
		To:     requestData.To,
		Time:   time.Now().UnixMilli(),
		Status: "delivered",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func SendMessageBulkHandler(w http.ResponseWriter, r *http.Request) {

	var requestData []model.SendMessageDataRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	var results []map[string]interface{}

	// Send the messages to each recipient
	for _, message := range requestData {
		result := map[string]interface{}{
			"to":      message.To,
			"type":    message.Type,
			"text":    message.Text,
			"caption": message.Caption,
			"url":     message.URL,
			"from":    message.From,
		}

		// Check if any required fields are missing
		if message.Type == "" || message.Text == "" || message.From == "" || message.To == "" {
			result["status"] = "failed"
			results = append(results, result)
			continue
		}

		// Validate 'from' number
		if !helpers.IsValidPhoneNumber(message.From) {
			result["status"] = "failed"
			results = append(results, result)
			continue
		}

		// Validate 'to' number
		if !helpers.IsValidPhoneNumber(message.To) {
			result["status"] = "failed"
			results = append(results, result)
			continue
		}

		// Validate if WhatsApp number is connected and give status failed if number not valid
		if !helpers.IsLoggedInByNumber(client, message.From) {
			result["status"] = "failed"
			results = append(results, result)
			continue
		}

		// Send the message if all checks pass
		if message.Type == "text" {
			err = helpers.SendMessageToPhoneNumber(client, message.To, message.Text)
			if err != nil {
				result["status"] = "failed"
			}
		} else {
			result["status"] = "failed"
		}

		// Add result to results slice
		results = append(results, result)
	}

	// Return the results as the response
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	data := make(map[string]map[string]interface{})
	for _, msg := range messages {
		timeStr := fmt.Sprintf("%d", msg.Time)
		// Remove @s.whatsapp.net suffix from msg.Chat
		chat := strings.TrimSuffix(msg.Chat, "@s.whatsapp.net")
		messageData := map[string]interface{}{
			"id":   msg.ID,
			"chat": chat,
			"time": msg.Time,
			"text": msg.Text,
		}
		data[timeStr] = messageData
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

func GetMessagesByIdHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	// Hapus akhiran @s.whatsapp.net dari id jika ada
	id = strings.TrimSuffix(id, "@s.whatsapp.net")

	data := make(map[string]map[string]interface{})
	for _, msg := range messages {
		// Periksa apakah nomor telepon terdapat dalam ID obrolan
		if strings.Contains(msg.Chat, id) {
			timeStr := fmt.Sprintf("%d", msg.Time)
			// Hapus akhiran @s.whatsapp.net dari msg.Chat
			chat := strings.TrimSuffix(msg.Chat, "@s.whatsapp.net")
			messageData := map[string]interface{}{
				"id":   msg.ID,
				"chat": chat,
				"time": msg.Time,
				"text": msg.Text,
			}
			data[timeStr] = messageData
		}
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

// CreateGroupHandler handles the creation of a new WhatsApp group
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Subject == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	if len(req.Participants) == 0 {
		http.Error(w, "At least one participant is required", http.StatusBadRequest)
		return
	}

	// Convert participant phone numbers to JID
	participants := make([]types.JID, len(req.Participants))
	for i, phone := range req.Participants {
		participantJID, err := types.ParseJID(phone + "@s.whatsapp.net")
		if err != nil {
			log.Printf("Error parsing JID for phone %s: %v", phone, err)
			http.Error(w, fmt.Sprintf("Invalid phone number: %s", phone), http.StatusBadRequest)
			return
		}
		participants[i] = participantJID
	}

	// Hypothetical field for group creation; replace with actual field names from the library
	reqCreateGroup := whatsmeow.ReqCreateGroup{
		// Try possible field names; here we assume 'Name' and 'Participants'
		Name:         req.Subject, // Hypothetical field
		Participants: participants,
	}

	groupResponse, err := client.CreateGroup(reqCreateGroup)
	if err != nil {
		log.Printf("Error creating group: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create group: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Group created successfully: %v", groupResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Group created successfully",
		"groupInfo": groupResponse,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the Logout method
	if err := client.Logout(); err != nil {
		http.Error(w, "Failed to log out user", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
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

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service is available")
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
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	version := "1.0.0" // Ganti dengan versi sistem yang sesuai
	response := model.VersionResponse{Version: version}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

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
	fmt.Printf("check device", deviceJIDs)
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
				"busy":    true,     // Replace with actual busy status if available
				"qrcode":  "",       // Replace with actual QR code if available
			}
			responseData = append(responseData, deviceData)
			continue // Continue to the next device
		}

		userInfo := userInfoMap[deviceJID]
		/*if !exists {
			http.Error(w, "User info not found for device", http.StatusNotFound)
			return
		}*/

		// Print userInfo for debugging
		fmt.Println("User Info:", userInfo)

		deviceData := map[string]interface{}{
			"id":      deviceJID.String(),
			"phone":   deviceJID.User,
			"status":  userInfo.Status,
			"process": "string", // Replace with actual process if available
			"busy":    false,    // Replace with actual busy status if available
			"qrcode":  "",       // Replace with actual QR code if available
		}

		responseData = append(responseData, deviceData)
	}

	response := map[string]interface{}{
		"data": responseData,
	}

	// Marshal the response into JSON and send it
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
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

func sendErrorResponse(w http.ResponseWriter, statusCode int, message, phoneNumber string) {
	errorResponse := response.ErrorResponseNumberPhone{
		StatusCode: statusCode,
		Error:      "Bad Request",
		Message:    "Device with phone: [" + phoneNumber + "] " + message,
	}

	// Convert ErrorResponse to JSON and send it as response
	w.WriteHeader(statusCode)
	jsonBytes, err := json.MarshalIndent(errorResponse, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// Handler untuk endpoint /api/grouplink
func GetGroupInviteLinkHandler(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")

	// Mendapatkan nilai reset dari parameter URL (atau dari body, sesuai kebutuhan)
	reset := r.URL.Query().Get("reset")

	// Lakukan validasi parameter jika diperlukan
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Konversi reset menjadi boolean
	resetBool := false
	if reset == "true" {
		resetBool = true
	}
	groupJID, err := types.ParseJID(groupID + "@g.us")

	inviteLink, err := client.GetGroupInviteLink(groupJID, resetBool)
	if err != nil {
		http.Error(w, "Failed to get group invite link", http.StatusInternalServerError)
		return
	}

	// Mengirimkan tautan undangan sebagai respons
	response := map[string]string{"group_invite_link": inviteLink}
	json.NewEncoder(w).Encode(response)
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
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
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
