package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	// Validasi platform
	if req.Platform == "whatsapp" {
		// (existing WhatsApp handling code)
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

		// Ambil thread_id dari request (jika ada)
		threadID := req.ThreadID

		// Kirim pesan menggunakan helper yang sesuai untuk Telegram dan sertakan thread_id
		if err := helpers.SendMessageToTelegram(combinedText); err != nil {
			helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidMessageType)
			return
		}

		// Kirim pesan ke thread lain jika perlu
		if err := helpers.SendMessageToTelegramThread(combinedText, threadID); err != nil {
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
