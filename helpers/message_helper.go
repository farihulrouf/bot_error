package helpers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
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

func UploadImageAndCreateMessage(client *whatsmeow.Client, imageBytes []byte, caption, mimeType string) (*waProto.ImageMessage, error) {
	// Unggah gambar
	resp, err := client.Upload(context.Background(), imageBytes, whatsmeow.MediaImage)
	if err != nil {
		return nil, fmt.Errorf("error uploading image: %v", err)
	}

	// Buat pesan gambar
	imageMsg := &waProto.ImageMessage{
		Caption:             proto.String(caption),
		Mimetype:            proto.String("image/jpeg"),
		ThumbnailDirectPath: &resp.DirectPath,
		ThumbnailSha256:     resp.FileSHA256,
		ThumbnailEncSha256:  resp.FileEncSHA256,
		//JpegThumbnail:       jpegBytes,

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}

	return imageMsg, nil
}

func MultipartFormFileHeaderToBytes(fileHeader *multipart.FileHeader) []byte {
	file, _ := fileHeader.Open()
	defer file.Close()

	fileBytes := make([]byte, fileHeader.Size)
	_, _ = file.Read(fileBytes)

	return fileBytes
}

func UploadDocAndCreateMessage(client *whatsmeow.Client, docBytes []byte, caption, mimeType string) (*waProto.DocumentMessage, error) {
	// Unggah Doc

	//client.Upload(context.Background(), data, whatsmeow.MediaDocument)
	//fileMimeType := http.DetectContentType(docBytes)

	resp, err := client.Upload(context.Background(), docBytes, whatsmeow.MediaDocument)
	if err != nil {
		return nil, fmt.Errorf("error uploading do: %v", err)
	}
	msg := &waProto.DocumentMessage{
		Url:      proto.String(resp.URL),
		Mimetype: proto.String(http.DetectContentType(docBytes)),
		//Title:         proto.String(resp.File.Filename),
		FileSha256: resp.FileSHA256,
		FileLength: proto.Uint64(resp.FileLength),
		MediaKey:   resp.MediaKey,
		//FileName:      proto.String(resp.File.Filename),
		FileEncSha256: resp.FileEncSHA256,
		DirectPath:    proto.String(resp.DirectPath),
		Caption:       proto.String(caption),
	}
	// Buat pesan Doc
	/*docMsg := &waProto.DocumentMessage{
		Caption: proto.String(caption),
		//Mimetype:            proto.String("image/jpeg"),
		Mimetype:            proto.String(http.DetectContentType(docBytes)),
		ThumbnailDirectPath: &resp.DirectPath,
		ThumbnailSha256:     resp.FileSHA256,
		ThumbnailEncSha256:  resp.FileEncSHA256,
		//JpegThumbnail:       jpegBytes,

		Url:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSha256: resp.FileEncSHA256,
		FileSha256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}*/
	//UploadDocAndCreateMessage
	return msg, nil
}

func UploadVideoAndCreateMessage(client *whatsmeow.Client, videoBytes []byte, caption, mimeType string) (*waProto.VideoMessage, error) {

	resp, err := client.Upload(context.Background(), videoBytes, whatsmeow.MediaVideo)
	if err != nil {
		return nil, fmt.Errorf("error uploading do: %v", err)
	}
	msg := &waProto.VideoMessage{
		Url:      proto.String(resp.URL),
		Mimetype: proto.String(http.DetectContentType(videoBytes)),
		//Title:         proto.String(resp.File.Filename),
		FileSha256: resp.FileSHA256,
		FileLength: proto.Uint64(resp.FileLength),
		MediaKey:   resp.MediaKey,
		//FileName:      proto.String(resp.File.Filename),
		FileEncSha256: resp.FileEncSHA256,
		DirectPath:    proto.String(resp.DirectPath),
		Caption:       proto.String(caption),
	}
	return msg, nil
}

func SendVideoToPhoneNumber(client *whatsmeow.Client, recipient string, videoMsg *waProto.VideoMessage) error {
	// Konversi recipient ke JID
	jid, err := types.ParseJID(recipient + "@s.whatsapp.net")
	if err != nil {
		return fmt.Errorf("invalid recipient JID: %v", err)
	}

	// Kirim pesan gambar
	_, err = client.SendMessage(context.Background(), jid, &waProto.Message{
		VideoMessage: videoMsg,
	})
	if err != nil {
		return fmt.Errorf("error sending image message: %v", err)
	}

	fmt.Printf("Sending image to phone number: %s\n", recipient)
	return nil
}

func SendImageToPhoneNumber(client *whatsmeow.Client, recipient string, imageMsg *waProto.ImageMessage) error {
	// Konversi recipient ke JID
	jid, err := types.ParseJID(recipient + "@s.whatsapp.net")
	if err != nil {
		return fmt.Errorf("invalid recipient JID: %v", err)
	}

	// Kirim pesan gambar
	_, err = client.SendMessage(context.Background(), jid, &waProto.Message{
		ImageMessage: imageMsg,
	})
	if err != nil {
		return fmt.Errorf("error sending image message: %v", err)
	}

	fmt.Printf("Sending image to phone number: %s\n", recipient)
	return nil
}

func SendDocToPhoneNumber(client *whatsmeow.Client, recipient string, docMsg *waProto.DocumentMessage) error {
	// Konversi recipient ke JID
	jid, err := types.ParseJID(recipient + "@s.whatsapp.net")
	if err != nil {
		return fmt.Errorf("invalid recipient JID: %v", err)
	}

	// Kirim pesan Doc
	_, err = client.SendMessage(context.Background(), jid, &waProto.Message{
		DocumentMessage: docMsg,
	})
	if err != nil {
		return fmt.Errorf("error sending doc message: %v", err)
	}

	fmt.Printf("Sending image to phone number: %s\n", recipient)
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
