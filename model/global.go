package model

import (
	"go.mau.fi/whatsmeow"
)

type DOConfig struct {
    Endpoint   string
	Bucket string
	Folder string
	AccessKey string
	SecretKey string
}

type CustomClient struct {
    User   int
	Phone  string
	ExpiredTime int64
	Webhook string
    Client *whatsmeow.Client
}

type ChangeWebhookRequest struct {
	Url string `json:"url"`
}

type PhoneRequest struct {
	Phone string `json:"phone"`
}

type PhoneParams struct {
	Phone string `json:"phone"`
}

type PhoneVerifyParams struct {
	Phone string `json:"phone"`
	PhoneRef string `json:"phoneref"`
	Ref string `json:"ref"`
}

type PhoneRefParams struct {
	Phone string `json:"phone"`
	Ref string `json:"ref"`
}

type PayloadWebhook struct {
	Section string `json:"section"`
	Data interface{} `json:"data"`
}

var Clients = make(map[string]CustomClient)
var SpaceConfig DOConfig
var DefaultWebhook string

