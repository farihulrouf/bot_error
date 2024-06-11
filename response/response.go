package response

type Message struct {
	ID   string `json:"id"`
	Chat string `json:"chat"`
	Time int64  `json:"time"`
	Text string `json:"text"`
	//Replies string `json:"text"`
}

type ErrorResponseNumberPhone struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}
