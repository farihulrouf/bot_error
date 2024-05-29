package model

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

// You can include more model definitions here if needed.
