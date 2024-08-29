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
    Client *whatsmeow.Client
}

var Clients = make(map[string]CustomClient)
var SpaceConfig DOConfig

