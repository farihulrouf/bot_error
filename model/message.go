// model/message.go

package model

// SendMessageGroupRequest defines the structure of the request to send a message to a group
type SendMessageDataRequest struct {
	To        string `json:"to"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	Caption   string `json:"caption"`
	URL       string `json:"url"`
	From      string `json:"from"`
	ImagePath string `json:"image_path"`
	MimeType  string `json:"mime_type"`
}
type SendMessageResponse struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Time   int64  `json:"time"`
	Status string `json:"status"`
}

type PhoneCodeRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type PhoneGroupRequest struct {
	Phone   string `json:"phone"`
	GroupID string `json:"group_id"`
}

// type Message struct {
// 	ID   string `json:"id"`
// 	Chat string `json:"chat"`
// 	Time int64  `json:"time"`
// 	Text string `json:"text"`
// 	/*
// 		Sender    string `json:"sender"`
// 		Message   string `json:"message"`
// 		Type      string `json:"type,omitempty"`       // Type of message (e.g., "text" or "media")
// 		MediaType string `json:"media_type,omitempty"` // Type of media (e.g., "image", "video", etc.)
// 		MediaURL  string `json:"media_url,omitempty"`  // URL of the media (if applicable)
// 	*/
// }

/*type CreateGroupRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
*/

type CreateGroupRequest struct {
	Subject      string   `json:"subject"` // Assuming 'subject' is the correct field for group name
	Participants []string `json:"participants"`
}

// LogoutRequest represents the payload for the /api/logout request
type LogoutRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type MessageData struct {
	ID        string `json:"id"`
	Time      int64  `json:"time"`
	FromMe    bool   `json:"fromMe"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	ChatType  string `json:"chatType"`
	ReplyID   string `json:"replyId"`
	Chat      string `json:"chat"`
	To        string `json:"to"`
	Name      string `json:"name"`
	From      string `json:"from"`
	Text      string `json:"text"`
	Caption   string `json:"caption"`
	URL       string `json:"url"`
	MimeType  string `json:"mimetype"`
	Thumbnail string `json:"thumbnail"`
}

type GetMessagesResponse struct {
	Data []MessageData `json:"data"`
}

type VersionResponse struct {
	Version string `json:"version"`
}

type GroupResponse struct {
	JID  string `json:"jid"`
	Name string `json:"name"`
}

type ReadReceipt struct {
	MessageID string `json:"message_id"`
	ReadBy    string `json:"read_by"`
	Time      int64  `json:"time"`
}

type Media struct {
	ID             	string `json:"id"`
	Type           	string `json:"type"`
	MimeType       	string `json:"mimeType"`
	Caption			string `json:"caption"`
	Thumbnail		[]byte `json:"thumbnail"`
	FileName       	string `json:"fileName"`
	FileContent	   	string `json:"fileContent"`
	FileLength		uint64 `json:"fileLength"`
	Seconds	   		uint32 `json:"seconds"`
	Name            string `json:"name"`
	Contact         string `json:"contact"`
	Poll			string `json:"poll"`
	Latitude		float64 `json:"latitude"`
	Longitude		float64 `json:"longitude"`
	Url 			string `json:"url"`
}

type Event struct {
	ID             string `json:"id"`
	Chat           string `json:"chat"`
	Time           int64  `json:"time"`
	SenderId	   string `json:"senderId"`
	SenderName     string `json:"senderName"`
	SenderAvatar   string `json:"senderAvatar"`
	IsGroup		   bool   `json:"isGroup"`
	Text           string `json:"text"`
	IsFromMe       bool   `json:"isFromMe"`
	MediaType	   string `json:"mediaType"`
	Type           string `json:"type"`
	ReplyToPost    string `json:"replyToPost"`
	ReplyToUser    string `json:"replyToUser"`
	IsViewOnce     bool   `json:"isViewOnce"`
	IsViewOnceV2   bool   `json:"isViewOnceV2"`
	IsViewOnceV2Extension   bool   `json:"isViewOnceV2Extension"`
	IsLottieSticker     	bool   `json:"isLottieSticker"`
	IsDocumentWithCaption   bool   `json:"isDocumentWithCaption"`
	Media 		   Media `json:"media"`
}

// WebhookRequest adalah struktur data untuk permintaan webhook
type WebhookRequest struct {
	URL string `json:"url"`
}

// WebhookResponse adalah struktur data untuk respons webhook
type WebhookResponse struct {
	Message string `json:"message"`
}

type WebhookPayload struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	To     string `json:"to"`
	Time   int64  `json:"time"`
	Status string `json:"status"`
}