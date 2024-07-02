package model

import (
	"database/sql"
	"fmt"
)

// Device struct represents the data model
type Device struct {
	JID            string `json:"jid"`
	RegistrationID string `json:"registration_id"`
	PushName       string `json:"pushname"`
}

// getDevicesFromDB function retrieves devices from the database
func GetDevicesFromDB(db *sql.DB) ([]Device, error) {
	query := `SELECT jid, registration_id, push_name FROM whatsmeow_device`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		err := rows.Scan(&device.JID, &device.RegistrationID, &device.PushName)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return devices, nil
}
