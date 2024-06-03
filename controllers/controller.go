package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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
	messages []model.Message
	mu       sync.Mutex
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
	groups, err := client.GetJoinedGroups()
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
}

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	groupJID, err := client.JoinGroupWithLink(req.InviteLink)
	if err != nil {
		log.Printf("Error joining group: %v", err)
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

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

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the message data
	var requestData model.SendMessageDataRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Send the message based on the recipient type
	if requestData.Type == "text" {
		err = helpers.SendMessageToPhoneNumber(client, requestData.To, requestData.Text)
	} else {
		http.Error(w, "Invalid message type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]string{"status": "Message sent successfully"}
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

	// Serialize messages to include both text and media messages
	var serializedMessages []interface{}
	for _, msg := range messages {
		serializedMsg := map[string]interface{}{
			"sender":  msg.Sender,
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

	// Filter messages by phone number
	var filteredMessages []model.Message
	for _, msg := range messages {
		if msg.Sender == id {
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
