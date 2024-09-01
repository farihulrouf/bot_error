package model

// User represents a user entity.
type Group struct {
	ID              string    `json:"id"`
	Name		    string    `json:"name"`
	Url		    	string    `json:"url"`
	OwnerID         string    `json:"owner_id"`
	IsIncognito     bool      `json:"is_incognito"`
	IsParent        bool      `json:"is_parent"`
	Avatar          string    `json:"avatar"`
	CreatedTime		int64	  `json:"created_time"`
	Members        []Member   `json:"members"`
}
