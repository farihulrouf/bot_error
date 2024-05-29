package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	waProto "go.mau.fi/whatsmeow/binary/proto"

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

func eventHandler(evt interface{}) {
	switch evt.(type) {
	default:
		fmt.Println("Unhandled event:", evt)
	}
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
	r := mux.NewRouter()
	r.HandleFunc("/api/groups", createGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups", getGroupsHandler).Methods("GET")
	r.HandleFunc("/api/groups/leave", leaveGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/join", JoinGroupHandler).Methods("POST")
	r.HandleFunc("/api/messages", sendMessageHandler).Methods("POST")

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

	// Join the group using the invite link
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
		RecipientType string `json:"recipient_type"`
		Recipient     string `json:"recipient"`
		Message       string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if recipient type is valid
	if requestData.RecipientType != "phone_number" && requestData.RecipientType != "group_id" {
		http.Error(w, "Invalid recipient type", http.StatusBadRequest)
		return
	}

	// Check if recipient and message are provided
	if requestData.Recipient == "" || requestData.Message == "" {
		http.Error(w, "Recipient and message are required", http.StatusBadRequest)
		return
	}

	// Send the message based on the recipient type
	//var err error
	switch requestData.RecipientType {
	case "phone_number":
		err = sendMessageToPhoneNumber(requestData.Recipient, requestData.Message)
	case "group_id":
		err = sendMessageToGroupID(requestData.Recipient, requestData.Message)
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
func PhoneNumberToJID(phoneNumber string) string {
	return phoneNumber + "@s.whatsapp.net"
}

// sendMessageToPhoneNumber sends a message to the specified phone number
func sendMessageToPhoneNumber(recipient, message string) error {
	// You would typically need to map the phone number to a JID
	// For simplicity, let's assume we have a way to do this
	//recipientJID := "1234567890@example.com" // Replace with the recipient's JID

	// Send the message using the global WhatsApp client

	// Send the message using the global WhatsApp client

	Jidsender, err := types.ParseJID(recipient + "@s.whatsapp.net")

	resp, err := client.SendMessage(context.Background(), Jidsender, &waProto.Message{
		Conversation: proto.String(message),
	})

	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	fmt.Printf("Sending message '%s' to group ID: %s\n", message, recipient, resp)

	return nil
}

// sendMessageToGroupID sends a message to the specified group ID
func sendMessageToGroupID(groupID, message string) error {
	// Replace this with the actual implementation using your messaging library
	fmt.Printf("Sending message '%s' to group ID: %s\n", message, groupID)
	// Example of sending the message:
	// result, err := messagingLibrary.SendMessageToGroupID(groupID, message)
	// if err != nil {
	//     return err
	// }
	// Handle result if needed
	return nil
}
