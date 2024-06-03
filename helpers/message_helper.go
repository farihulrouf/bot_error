package helpers

import (
	"context"
	"fmt"

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

func SendMessage(client *whatsmeow.Client, jid types.JID, req model.SendMessageGroupRequest) error {
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
