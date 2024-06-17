package response

type Message struct {
	ID             string `json:"id"`
	Chat           string `json:"chat"`
	Time           int64  `json:"time"`
	Text           string `json:"text"`
	Group          bool   `json:"group,omitempty"`
	IsFromMe       bool   `json:"isfromme,omitempty"`
	CommentMessage string `json:"comment"`
	Tipe           string `json:"tipe"`
	IsDocument     bool   `json:"isdocument,omitempty"`
	Mediatipe      string `json:"mediatipe"`
	Caption        string `json:"caption"`
	VideoMessage   string `json:"videomessage"`
	DocMessage     string `json:"docmessage"`
	MimeTipe       string `json:"mimetipe"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	From           string `json:"from"`
	To             string `json:"to"`
	Thumbnail      string `json:"thumbnail"`
	Thumbnailvideo string `json:"thumbnailvideo"`
	Thumbnaildoc   string `json:"thumbnaildoc"`

	//Replies string `json:"text"`
}

type ErrorResponseNumberPhone struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

type GroupResponse struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
	Admins      []string `json:"admins"`
	Time        int64    `json:"time"`
	Pinned      bool     `json:"pinned"`
	UnreadCount int      `json:"unreadCount"`
}

//

type MessageFilter struct {
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

type Response struct {
	Data []MessageFilter `json:"data"`
}
