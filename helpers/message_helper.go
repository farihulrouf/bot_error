package helpers

import (
	"context"
	"fmt"
	"strings"
	"time"

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
