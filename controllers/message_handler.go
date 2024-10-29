package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"wagobot.com/errors"
	"wagobot.com/helpers"
	"wagobot.com/model"
)

func SendMessageGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.SendMessageDataRequest

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidRequestPayload)
		return
	}

	if req.Platform == "whatsapp" {
		var value_client = model.Clients["device1"].Client
		matchFound := false

		// Validasi field Message
		if req.Message == "" {
			helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrMissingRequiredFields)
			return
		}

		// Membuat text gabungan
		combinedText := fmt.Sprintf(
			"Name: %s\nIP: %s\nTime: %d\nStatus: %s\nMessage: %s",
			req.Name, req.IP, req.Time, req.Status, req.Message,
		)

		// Set text pada request
		req.Text = combinedText

		// Set from, to, and type secara otomatis
		req.From = "6285280933757"         // Set nomor pengirim
		req.To = "120363314404357759@g.us" // Set grup tujuan
		req.Type = "text"                  // Set tipe pesan ke "text"

		// Convert to JID
		jid, err := helpers.ConvertToJID(req.To)
		if err != nil {
			helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidRecipient)
			return
		}

		// Cek klien yang cocok
		for key := range model.Clients {
			whoami := model.Clients[key].Client.Store.ID.String()
			parts := strings.Split(whoami, ":")
			if req.From == parts[0] {
				value_client = model.Clients[key].Client
				matchFound = true
				break
			}
		}

		if !matchFound {
			helpers.SendErrorResponse(w, http.StatusBadRequest, "No matching number found for requestData.From")
			return
		}

		// Kirim pesan menggunakan value_client dan jid
		if err := helpers.SendMessage(value_client, jid, req); err != nil {
			helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidMessageType)
			return
		}

		// Respons dengan sukses
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Message sent successfully"})
	} else if req.Platform == "telegram" {
		// Logic for handling Telegram messages
		if req.Message == "" {
			helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrMissingRequiredFields)
			return
		}

		// Membuat text gabungan untuk Telegram
		combinedText := fmt.Sprintf(
			"Name: %s\nIP: %s\nTime: %d\nStatus: %s\nMessage: %s",
			req.Name, req.IP, req.Time, req.Status, req.Message,
		)

		// Kirim pesan menggunakan helper yang sesuai untuk Telegram
		if err := helpers.SendMessageToTelegram(combinedText); err != nil {
			helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidMessageType)
			return
		}

		// Respons dengan sukses
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Message sent successfully"})
	} else {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Unsupported platform")
	}
}
