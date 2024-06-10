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
)

var client *whatsmeow.Client

func SetClient(c *whatsmeow.Client) {
	client = c
}

var (
	messages   []model.Message
	mu         sync.Mutex
	webhookURL string
)

func sendToAPI(sender string, message string) {
	mu.Lock()
	messages = append(messages, model.Message{Sender: sender, Message: message})
	mu.Unlock()
}
func EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if !v.Info.IsFromMe && v.Message.GetConversation() != "" {
			fmt.Println("PESAN DITERIMA!", v.Message.GetConversation())
			sendToAPI(v.Info.Sender.String(), v.Message.GetConversation())
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

	var filteredGroups []model.GroupResponse
	for _, group := range groups {
		for _, member := range group.Participants {
			fmt.Println("Participatan Group List", member)
			if member.JID.User == phone {
				fmt.Println("Member Jid User", member.JID.User)
				filteredGroups = append(filteredGroups, model.GroupResponse{JID: group.JID.String(), Name: group.Name})
				break
			}
		}
	}

	response, err := json.Marshal(filteredGroups)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	/*groups, err := client.GetJoinedGroups()
		if err != nil {
			http.Error(w, "Failed to get groups", http.StatusInternalServerError)
			return
		}

		groupList := make([]map[string]interface{}, 0, len(groups))
		for _, group := range groups {
			fmt.Printf("Group ID: %s, Name: %s\n", group.JID.String(), group.Name)
			groupInfo := map[string]interface{}{
				"JID":  group.JID.String(),
				"Name": group.Name,
			}
			groupList = append(groupList, groupInfo)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groupList)
	}

	func convertToJID(to string) (types.JID, error) {
		var jid types.JID
		var err error

		if strings.Contains(to, "-") {
			// Assuming it's a Group ID
			jid, err = types.ParseJID(to + "@g.us")
		} else {
			// Assuming it's a phone number
			jid, err = types.ParseJID(to + "@s.whatsapp.net")
		}

		if err != nil {
			return types.JID{}, fmt.Errorf("invalid JID: %v", err)
		}

		return jid, nil
	*/
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
		http.Error(w, "Recipient number is not connected to WhatsApp", http.StatusBadRequest)
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
	// Parse request body to get the message data
	var requestData []model.SendMessageDataRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Send the messages to each recipient
	var successCount int
	for _, message := range requestData {
		// Check if any required fields are missing
		if message.Type == "" || message.Text == "" || message.From == "" || message.To == "" {
			fmt.Printf("Missing required fields for recipient %s\n", message.To)
			continue // Skip sending the message for this recipient
		}

		// Validate 'from' number
		if !helpers.IsValidPhoneNumber(message.From) {
			fmt.Printf("Invalid 'from' number for recipient %s: %s\n", message.To, message.From)
			continue // Skip sending the message for this recipient
		}

		// Validate 'to' number
		if !helpers.IsValidPhoneNumber(message.To) {
			fmt.Printf("Invalid 'to' number for recipient %s: %s\n", message.To, message.To)
			continue // Skip sending the message for this recipient
		}

		// Send the message if all checks pass
		if message.Type == "text" {
			err = helpers.SendMessageToPhoneNumber(client, message.To, message.Text)
			if err != nil {
				fmt.Printf("Failed to send message to %s: %v\n", message.To, err)
			} else {
				successCount++
			}
		} else {
			fmt.Printf("Invalid message type for recipient %s\n", message.To)
		}
	}

	// Return success response
	response := map[string]interface{}{
		"status":           "Bulk message sent",
		"success_count":    successCount,
		"total_recipients": len(requestData),
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

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Lock messages to ensure thread safety
	mu.Lock()
	defer mu.Unlock()

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	var serializedMessages []interface{}
	for _, msg := range messages {
		//sender := msg.Sender
		sender := strings.Split(msg.Sender, "@")[0]
		//sender_phone := sender[0]
		serializedMsg := map[string]interface{}{
			"sender":  sender,
			"type":    msg.Type,
			"message": msg.Message,
		}

		// Add additional fields for media messages if present
		if msg.Type == "media" {
			serializedMsg["media_type"] = msg.MediaType
			serializedMsg["media_url"] = msg.MediaURL
		}

		serializedMessages = append(serializedMessages, serializedMsg)
	}

	// Encode messages array to JSON and send response
	if err := json.NewEncoder(w).Encode(serializedMessages); err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
		return
	}
}

func GetMessagesByIdHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Extract id from the request URL path
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]
	//fmt.Println("check number phone", id)
	// Filter messages by phone number
	var filteredMessages []model.Message
	for _, msg := range messages {
		if msg.Sender == id+"@s.whatsapp.net" {
			filteredMessages = append(filteredMessages, msg)
		}
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Encode filtered messages array to JSON and send response
	if err := json.NewEncoder(w).Encode(filteredMessages); err != nil {
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
