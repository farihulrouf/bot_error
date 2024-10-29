package helpers

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"wagobot.com/model"
)

func SendMessageToTelegram(message string) error {
	// Ambil chat ID dari variabel lingkungan
	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if chatID == "" {
		return fmt.Errorf("TELEGRAM_CHAT_ID not set")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}

	apiURL := os.Getenv("TELEGRAM_API_URL")
	if apiURL == "" {
		return fmt.Errorf("TELEGRAM_API_URL not set")
	}

	url := fmt.Sprintf("%s%s/sendMessage", apiURL, botToken)

	body, err := json.Marshal(map[string]string{
		"chat_id": chatID,
		"text":    message,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message: %s", resp.Status)
	}

	return nil
}

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
	imageMsg := &waE2E.ImageMessage{
		Caption:       proto.String(caption),
		Mimetype:      proto.String("image/png"),
		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
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
	/*msg := &waProto.DocumentMessage{
		URL:      proto.String(resp.URL),
		Mimetype: proto.String(http.DetectContentType(docBytes)),
		//Title:         proto.String(resp.File.Filename),
		FileSHA256: resp.FileSHA256,
		FileLength: proto.Uint64(resp.FileLength),
		MediaKey:   resp.MediaKey,
		//FileName:      proto.String(resp.File.Filename),
		FileEncSHA256: resp.FileEncSHA256,
		DirectPath:    proto.String(resp.DirectPath),
		Caption:       proto.String(caption),
	}*/

	msg := &waE2E.DocumentMessage{ // Nama file dari URL atau SHA-256 sebagai title		// Nama file atau SHA-256 sebagai title
		Mimetype:      proto.String(http.DetectContentType(docBytes)), // MIME type dokumen, misalnya "application/pdf" untuk PDF
		URL:           &resp.URL,                                      // URL tempat file diupload
		DirectPath:    &resp.DirectPath,                               // Direct path untuk file
		MediaKey:      resp.MediaKey,                                  // Kunci enkripsi untuk media
		FileEncSHA256: resp.FileEncSHA256,                             // SHA-256 dari file yang dienkripsi
		FileSHA256:    resp.FileSHA256,                                // SHA-256 dari file asli
		FileLength:    &resp.FileLength,                               // Ukuran file dalam byte
		Caption:       proto.String(caption),                          // Nama file yang akan didownload
	}

	//UploadDocAndCreateMessage
	return msg, nil
}

func UploadVideoAndCreateMessage(client *whatsmeow.Client, videoBytes []byte, caption, mimeType string) (*waProto.VideoMessage, error) {

	resp, err := client.Upload(context.Background(), videoBytes, whatsmeow.MediaVideo)
	if err != nil {
		return nil, fmt.Errorf("error uploading do: %v", err)
	}
	msg := &waE2E.VideoMessage{
		URL:      proto.String(resp.URL),
		Mimetype: proto.String(http.DetectContentType(videoBytes)),
		//Title:         proto.String(resp.File.Filename),
		FileSHA256: resp.FileSHA256,
		FileLength: proto.Uint64(resp.FileLength),
		MediaKey:   resp.MediaKey,
		//FileName:      proto.String(resp.File.Filename),
		FileEncSHA256: resp.FileEncSHA256,
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
func IsValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	// Periksa bahwa skema (scheme) URL adalah http atau https
	return u.Scheme == "http" || u.Scheme == "https"
}

// Fungsi untuk mendeteksi ekstensi file dan jenis file
func DetectFileType(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "Error parsing URL"
	}

	// Mendapatkan path dari URL
	path := parsedURL.Path
	// Mendapatkan ekstensi file
	ext := strings.ToLower(strings.TrimPrefix(path[strings.LastIndex(path, "."):], "."))

	// Menentukan jenis file berdasarkan ekstensi
	switch {
	case ext == "jpg" || ext == "jpeg" || ext == "png" || ext == "gif" || ext == "bmp":
		return "Image"
	case ext == "mp4" || ext == "mkv" || ext == "avi" || ext == "mov" || ext == "webm":
		return "Video"
	case ext == "pdf" || ext == "doc" || ext == "docx" || ext == "xls" || ext == "xlsx" || ext == "ppt" || ext == "pptx":
		return "Document"
	case ext == "mp3" || ext == "wav" || ext == "flac" || ext == "ogg":
		return "Audio"
	default:
		return "Unknown type"
	}
}

func DownloadFile(urlStr string) ([]byte, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DetectFileTypeByContentType(urlStr string) (string, error) {
	resp, err := http.Head(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix(contentType, "image/"):
		// Deteksi berbagai jenis gambar
		switch {
		case contentType == "image/jpeg":
			return "Image: JPEG", nil
		case contentType == "image/png":
			return "Image: PNG", nil
		case contentType == "image/gif":
			return "Image: GIF", nil
		case contentType == "image/webp":
			return "Image: WEBP", nil
		case contentType == "image/bmp":
			return "Image: BMP", nil
		case contentType == "image/svg+xml":
			return "Image: SVG", nil
		case contentType == "image/tiff":
			return "Image: TIFF", nil
		default:
			return "Image: Unknown format", nil
		}
	case strings.HasPrefix(contentType, "video/"):
		// Deteksi berbagai jenis video
		switch {
		case contentType == "video/mp4":
			return "Video: MP4", nil
		case contentType == "video/avi":
			return "Video: AVI", nil
		case contentType == "video/mkv":
			return "Video: MKV", nil
		case contentType == "video/webm":
			return "Video: WEBM", nil
		case contentType == "video/quicktime":
			return "Video: QuickTime", nil
		default:
			return "Video: Unknown format", nil
		}
	case strings.HasPrefix(contentType, "application/pdf"):
		return "Document: PDF", nil
	case strings.HasPrefix(contentType, "application/vnd.openxmlformats-officedocument/wordprocessingml.document"):
		return "Document: DOCX", nil
	case strings.HasPrefix(contentType, "application/msword"):
		return "Document: DOC", nil
	case strings.HasPrefix(contentType, "application/vnd.openxmlformats-officedocument/spreadsheetml.sheet"):
		return "Document: XLSX", nil
	case strings.HasPrefix(contentType, "application/vnd.ms-excel"):
		return "Document: XLS", nil
	case strings.HasPrefix(contentType, "application/vnd.openxmlformats-officedocument/presentationml.presentation"):
		return "Document: PPTX", nil
	case strings.HasPrefix(contentType, "application/vnd.ms-powerpoint"):
		return "Document: PPT", nil
	case strings.HasPrefix(contentType, "audio/"):
		// Deteksi berbagai jenis audio
		switch {
		case contentType == "audio/mpeg":
			return "Audio: MP3", nil
		case contentType == "audio/wav":
			return "Audio: WAV", nil
		case contentType == "audio/ogg":
			return "Audio: OGG", nil
		case contentType == "audio/flac":
			return "Audio: FLAC", nil
		default:
			return "Audio: Unknown format", nil
		}
	default:
		return "Unknown type", nil
	}
}

// Fungsi untuk memeriksa apakah string adalah Base64 yang valid
func IsBase64(str string) bool {
	_, err := base64.StdEncoding.DecodeString(str)
	return err == nil
}

// Fungsi untuk mengonversi Base64 ke byte array dan menghitung hash dari data
func Base64ToHash(base64Str string) ([]byte, []byte, error) {
	// Dekode string Base64 menjadi byte array
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, nil, err
	}

	// Hitung hash SHA-256
	hashSHA256 := sha256.New()
	hashSHA256.Write(data)
	hashSHA256Bytes := hashSHA256.Sum(nil)

	// Hitung hash MD5
	hashMD5 := md5.New()
	hashMD5.Write(data)
	hashMD5Bytes := hashMD5.Sum(nil)

	return hashSHA256Bytes, hashMD5Bytes, nil
}
func encryptMD5(chatId string) string {
	hash := md5.New()
	hash.Write([]byte(chatId))
	return hex.EncodeToString(hash.Sum(nil))
}

func ConvertToLettersDetailed(number string, chatId string, isGroup bool) string {
	encryptedChatId := encryptMD5(chatId)
	var result strings.Builder
	if isGroup {
		result.WriteString(encryptedChatId + "G")
	} else {
		result.WriteString(encryptedChatId + "U")

	}

	return result.String()
}
