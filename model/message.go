// model/message.go

package model

// SendMessageGroupRequest defines the structure of the request to send a message to a group
type SendMessageDataRequest struct {
	To      string `json:"to"`
	Type    string `json:"type"`
	Text    string `json:"text"`
	Caption string `json:"caption"`
	URL     string `json:"url"`
	From    string `json:"from"`
}

type JoinGroupRequest struct {
	InviteLink string `json:"invite_link"`
}

type LeaveGroupRequest struct {
	GroupID string `json:"group_id"`
	Phone   string `json:"phone"`
}

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Type      string `json:"type,omitempty"`       // Type of message (e.g., "text" or "media")
	MediaType string `json:"media_type,omitempty"` // Type of media (e.g., "image", "video", etc.)
	MediaURL  string `json:"media_url,omitempty"`  // URL of the media (if applicable)
}

/*type CreateGroupRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
*/

type CreateGroupRequest struct {
	Subject      string   `json:"subject"` // Assuming 'subject' is the correct field for group name
	Participants []string `json:"participants"`
}
