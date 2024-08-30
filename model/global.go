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
	ExpiredTime int64
	Webhook string
    Client *whatsmeow.Client
}

type ChangeWebhookRequest struct {
	Url string `json:"url"`
}

var Clients = make(map[string]CustomClient)
var SpaceConfig DOConfig
var DefaultWebhook string

