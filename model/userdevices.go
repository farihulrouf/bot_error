package model

import "strings"

// User represents a user entity.
type UserDevice struct {
	ID              int       `json:"id"`
	UserId        	int       `json:"user_id"`
	DeviceJid       string    `json:"device_jid"`
}

func GetPhoneNumber(jid string) string {
	parts := strings.Split(jid, "@")
    parts = strings.Split(parts[0], ":")
    return parts[0]
}

func GetPhoneFromDeviceScan(devid string) string {
	// format DEVICE-{phone}-{ref}
	parts := strings.Split(devid, "-")
    return parts[1]
}