package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	//"github.com/tulir/whatsmeow/binary/proto"

	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

type CreateGroupRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type LeaveGroupRequest struct {
	GroupID string `json:"group_id"`
	Phone   string `json:"phone"`
}

type GroupMessage struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type JoinGroupRequest struct {
	InviteLink string `json:"invite_link"`
}

type Message struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

/*func eventHandler(evt interface{}) {
	switch evt.(type) {
	default:
		fmt.Println("Unhandled event:", evt)
	}
}
*/

var (
	messages []Message
	mu       sync.Mutex
)

func sendToAPI(sender string, message string) {
	mu.Lock()
	messages = append(messages, Message{Sender: sender, Message: message})
	mu.Unlock()
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if !v.Info.IsFromMe && v.Message.GetConversation() != "" {
			fmt.Println("PESAN DITERIMA!", v.Message.GetConversation())
			sendToAPI(v.Info.Sender.String(), v.Message.GetConversation())
		}
	}
}

type SendMessageGroupRequest struct {
	GroupID string `json:"group_id"`
	Message string `json:"message"`
}

func addSampleMessage() {
	sendToAPI("6282333899903", "Hello, this is a sample message!")
}

func main() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:wasopingi.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	ScanQrCode(client)

	// Setup router
	// Add a sample message for testing
	addSampleMessage()
	r := mux.NewRouter()
	r.HandleFunc("/api/groups", createGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups", getGroupsHandler).Methods("GET")
	r.HandleFunc("/api/groups/messages", sendMessageGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/leave", leaveGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/join", JoinGroupHandler).Methods("POST")
	r.HandleFunc("/api/messages", sendMessageHandler).Methods("POST")
	r.HandleFunc("/api/result", getMessages).Methods("GET")
	r.HandleFunc("/api/result/{id}", getMessagesByPhoneNumber).Methods("GET")
	r.HandleFunc("/api/messages/bulk", sendMessageBulkHandler).Methods("POST")
	//r.HandleFunc("/api/messages/{phone}", getMessagesByPhoneNumber).Methods("GET")
	//r.HandleFunc("/api/messages/{phone}", getMessagesByPhoneNumber).Methods("GET")

	//sendMessageGroupHandler
	//r.HandleFunc("/api/messages", GetAllMessagesHandler).Methods("GET")

	// Start server
	go func() {
		log.Println("Server running on port 8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	participantJID, err := types.ParseJID(req.Phone + "@s.whatsapp.net")
	if err != nil {
		log.Printf("Error parsing JID: %v", err)
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	reqCreateGroup := whatsmeow.ReqCreateGroup{
		Participants: []types.JID{participantJID},
	}

	groupResponse, err := client.CreateGroup(reqCreateGroup)
	if err != nil {
		log.Printf("Error creating group: %v", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	log.Printf("Group created successfully: %v", groupResponse)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group created successfully"})
}

func getGroupsHandler(w http.ResponseWriter, r *http.Request) {
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

func leaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req LeaveGroupRequest
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

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req JoinGroupRequest
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

// sendMessageHandler handles sending a message to either a phone number or a group ID
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the message data
	var requestData struct {
		To      string `json:"to"`
		Type    string `json:"type"`
		Text    string `json:"text"`
		Caption string `json:"caption"`
		URL     string `json:"url"`
		From    string `json:"from"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Send the message based on the recipient type
	if requestData.Type == "text" {
		err = sendMessageToPhoneNumber(requestData.To, requestData.Text)
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

// sendMessageToPhoneNumber sends a message to the specified phone number
func sendMessageToPhoneNumber(recipient, message string) error {
	// Convert recipient to JID
	jid, err := types.ParseJID(recipient + "@s.whatsapp.net")
	if err != nil {
		return fmt.Errorf("invalid recipient JID: %v", err)
	}

	// Create the message
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	// Send the message
	_, err = client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	fmt.Printf("Sending message '%s' to phone number: %s\n", message, recipient)
	return nil
}

func sendMessageGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req SendMessageGroupRequest

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Send the message
	if err := sendMessageToGroupID(req.GroupID, req.Message); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Message sent to group ID: %s", req.GroupID)
}

// sendMessageToGroupID sends a message to the specified group ID
func sendMessageToGroupID(groupID, message string) error {
	// Convert groupID to JID
	fmt.Println("Chek groupid", groupID)
	jid, err := types.ParseJID(groupID + "@g.us")
	if err != nil {
		return fmt.Errorf("invalid group JID: %v", err)
	}

	// Create the message
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	// Send the message
	_, err = client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	fmt.Printf("Sending message '%s' to group ID: %s\n", message, groupID)
	return nil
}

func sendMessageBulkHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the message data
	var requestData []struct {
		To      string `json:"to"`
		Type    string `json:"type"`
		Text    string `json:"text"`
		Caption string `json:"caption"`
		URL     string `json:"url"`
		From    string `json:"from"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Send the messages to each recipient
	var successCount int
	for _, message := range requestData {
		if message.Type == "text" {
			err = sendMessageToPhoneNumber(message.To, message.Text)
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

func getMessages(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// NormalizePhoneNumber normalizes the phone number format
func normalizePhoneNumber(phone string) string {
	// Add any necessary normalization logic here
	return phone
}

func getMessagesByPhoneNumber(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	phone := params["phone"]
    
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Checking messages for phone:", phone)
	var filteredMessages []Message
	for _, message := range messages {
		fmt.Println("Checking message from sender:", message.Sender)
		if message.Sender == phone {
			filteredMessages = append(filteredMessages, message)
		}
	}

	if len(filteredMessages) == 0 {
		http.Error(w, "No messages found for this phone number", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredMessages)
}
