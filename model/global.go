package model

import (
	"go.mau.fi/whatsmeow"
)

type CustomClient struct {
    User   int
	ExpiredTime int64
    Client *whatsmeow.Client
}

var Clients = make(map[string]CustomClient)

