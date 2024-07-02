package helpers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"wagobot.com/model"
)

func SendMessageToPhoneNumber(client *whatsmeow.Client, recipient, message string) error {
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
func SendMessage(client *whatsmeow.Client, jid types.JID, req model.SendMessageDataRequest) error {
	// Create the message based on the type
	var msg *waProto.Message
	switch req.Type {
	case "text":
		msg = &waProto.Message{
			Conversation: proto.String(req.Text),
		}
		// Add more cases for different message types as needed
	default:
		return fmt.Errorf("unsupported message type: %s", req.Type)
	}

	// Send the message
	_, err := client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	fmt.Printf("Sending message '%s' to %s from %s\n", req.Text, jid.String(), req.From)
	return nil
}

func ConvertToJID(to string) (types.JID, error) {
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

func SendMessageToGroup(client *whatsmeow.Client, groupJID types.JID, message string) error {
	// Create the message
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	// Send the message to the group
	_, err := client.SendMessage(context.Background(), groupJID, msg)
	if err != nil {
		return fmt.Errorf("error sending message to group %s: %v", groupJID, err)
	}

	fmt.Printf("Sending message '%s' to group %s\n", message, groupJID.String())
	return nil
}

// GetAllMessagesByPhoneNumberOrGroupID gets all messages by phone number or group ID
func GetAllMessagesByPhoneNumberOrGroupID(client *whatsmeow.Client, identifier string) ([]model.MessageData, error) {
	// Simulate fetching messages
	// Replace this with actual logic to get messages from WhatsApp
	messages := []model.MessageData{
		{
			ID:        "1609773514305",
			Time:      time.Now().UnixMilli(),
			FromMe:    true,
			Type:      "text",
			Status:    "delivered",
			ChatType:  "user",
			ReplyID:   "1609773514305",
			Chat:      identifier,
			To:        "6281234567890",
			Name:      "string",
			From:      identifier,
			Text:      "Test from MaxChat",
			Caption:   "Caption test",
			URL:       "https://www.fnordware.com/superpng/pnggrad16rgb.png",
			MimeType:  "string",
			Thumbnail: "string",
		},
	}

	return messages, nil
}

func IsValidPhoneNumber(phoneNumber string) bool {
	// Parse the phone number using "ZZ" as the default region (which allows parsing international numbers)
	num, err := phonenumbers.Parse(phoneNumber, "ID")
	if err != nil {
		fmt.Println("Error parsing phone number:", err)
		return false
	}
	// Check if the phone number is valid
	return phonenumbers.IsValidNumber(num)
}

func UpdateUserInfo(values map[string]string, field string, value string) map[string]string {
	log.Debug().Str("field", field).Str("value", value).Msg("User info updated")
	values[field] = value
	return values
}

func IsLoggedInByNumber(client *whatsmeow.Client, phoneNumber string) bool {
	// Memeriksa status login berdasarkan nomor telepon menggunakan IsLoggedIn
	//fmt.Println("check numberphone", phoneNumber)
	if !client.IsLoggedIn() {
		return false
	}

	deviceStore := client.Store.ID
	if deviceStore == nil {
		fmt.Println("Device store is nil")
		return false
	}

	// Convert the deviceStore ID to a proper JID
	selfJID := types.NewJID(deviceStore.User, types.DefaultUserServer)
	userInfoMap, err := client.GetUserInfo([]types.JID{selfJID})
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return false
	}
	userInfo := userInfoMap[selfJID]
	fmt.Printf("check device : %v\n", userInfo)

	// Dapatkan nomor telepon dari userInfoMap
	nomorTeleponDitemukan := false
	for _, userInfo := range userInfoMap {

		for _, device := range userInfo.Devices {
			if device.User == phoneNumber {
				fmt.Printf("Nomor Telepon: %s\n", phoneNumber)
				nomorTeleponDitemukan = true
				break
			}
		}

	}

	if !nomorTeleponDitemukan {
		fmt.Println("Nomor Telepon tidak ditemukan dalam userInfoMap")
		return false
	}

	// Nomor Telepon ditemukan
	return true
}
