// model/message.go

package model

// SendMessageGroupRequest defines the structure of the request to send a message to a group
type SendMessageGroupRequest struct {
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
	Sender  string `json:"sender"`
	Message string `json:"message"`
}
