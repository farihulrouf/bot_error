package model

// User represents a user entity.
type Member struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	Name           	string    `json:"name"`
	Phone           string    `json:"phone"`
	Avatar          string    `json:"avatar"`
	GroupID       	string 	  `json:"group_id"`
	IsAdmin         bool      `json:"is_admin"`
	IsSuperAdmin    bool      `json:"is_super_admin"`
}
