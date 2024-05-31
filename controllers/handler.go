package main

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// ScanQrCode memindai dan mencetak kode QR untuk login.
func ScanQrCode(client *whatsmeow.Client) {
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Print each line of the QR code
				for _, line := range evt.Code {
					fmt.Println("print line Qr", line)
				}
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		fmt.Println("Already logged in.")
		err := client.Connect()
		if err != nil {
			panic(err)
		}
	}
}

// terinma semua pasan memproses pesan masuk dan mencetak isinya.
func ReceiveAllMessages(v *events.Message) {
	if !v.Info.IsFromMe {
		sender := v.Info.Sender.String()
		switch {
		case v.Message.GetConversation() != "":
			fmt.Printf("Text message from %s: %s\n", sender, v.Message.GetConversation())
		case v.Message.GetImageMessage() != nil:
			fmt.Printf("Image message from %s: %s\n", sender, v.Message.GetImageMessage().GetCaption())
		case v.Message.GetVideoMessage() != nil:
			fmt.Printf("Video message from %s: %s\n", sender, v.Message.GetVideoMessage().GetCaption())
		case v.Message.GetAudioMessage() != nil:
			fmt.Printf("Audio message from %s\n", sender)
		case v.Message.GetDocumentMessage() != nil:
			fmt.Printf("Document message from %s: %s\n", sender, v.Message.GetDocumentMessage().GetTitle())
		default:
			fmt.Printf("Other type of message from %s\n", sender)
		}
	}
}

func ManageGroups() {
	// Terapkan logika manajemen grup di sini
	fmt.Println("Manajemen login di sini")
}

// // HandleGroupChangeEvent menangani event yang berhubungan dengan perubahan grup.
func HandleGroupChangeEvent(evt interface{}) {
	// // Untuk saat ini, cetak saja event tersebut untuk menunjukkan bahwa event tersebut tidak tertangani
	fmt.Println("Unhandled event:", evt)
}
