package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	//"go.mau.fi/whatsmeow"
	"wagobot.com/errors"
	"wagobot.com/helpers"
	"wagobot.com/model"
	"wagobot.com/response"
)

// Pastikan client diimpor dari file yang sesuai
//var client *whatsmeow.Client

func SendMessageGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.SendMessageDataRequest

	// Decode the JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidRequestPayload)
		return
	}

	// Validate the request data
	if req.Type == "" || req.Text == "" {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrMissingRequiredFields)
		return
	}

	// Convert to JID
	jid, err := helpers.ConvertToJID(req.To)
	if err != nil {
		http.Error(w, fmt.Sprintf(errors.ErrInvalidRecipient, ":%v", err), http.StatusBadRequest)
		return
	}

	// Send the message
	if err := helpers.SendMessage(client, jid, req); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf(errors.ErrInvalidMessageType, ":%v", err))
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Message sent to: %s", req.To)
}

// SendMessageHandler handles sending messages.

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the message data
	//var isAdmin bool
	//adminGroupJIDs := make([]string, 0)

	var requestData model.SendMessageDataRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Check if any required field is missing
	if requestData.To == "" || requestData.Type == "" || requestData.Text == "" || requestData.From == "" {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Missing required fields: 'to', 'type', 'text', or 'from'")
		return
	}

	// Validate phone numbers
	if !helpers.IsValidPhoneNumber(requestData.To) {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidPhoneNumberTo)
		return
	}
	if !helpers.IsValidPhoneNumber(requestData.From) {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidPhoneNumberSender)
		return
	}
	//Check numerphone is login or not
	if !helpers.IsLoggedInByNumber(client, requestData.From) {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrNotReadyOrNotAvailable)
		return
	}
	// Simpan daftar JID grup yang merupakan admin

	// Periksa setiap grup untuk memeriksa apakah pengguna adalah admin
	/* tidak perlu mengirim pesan ke group hanya private / only private
	groups, err := client.GetJoinedGroups()
	if err != nil {
		http.Error(w, "Failed to fetch joined groups", http.StatusInternalServerError)
		return
	}

	///looping for message to group
	for _, group := range groups {
		for _, participant := range group.Participants {
			if participant.JID.User == requestData.To && participant.IsAdmin {
				// Jika nomor tersebut adalah admin, tambahkan JID grup ke dalam slice
				adminGroupJIDs = append(adminGroupJIDs, group.JID.String())

				//fmt.Println("check admin", adminGroupJIDs)
				// Keluar dari loop inner karena sudah ditentukan bahwa nomor tersebut adalah admin dalam grup ini
				break
			}
		}
	}

	fmt.Println("nomer jid", adminGroupJIDs)
	for _, groupJID := range adminGroupJIDs {
		// Convert the string JID to types.JID
		//jid := types.JID(groupJID)
		parts := strings.Split(groupJID, "@")

		// Extract user and server parts
		user := parts[0]
		server := parts[1]

		// Convert the user and server parts to types.JID
		jid := types.NewJID(user, server)

		// Call the SendMessageToGroup function
		err := helpers.SendMessageToGroup(client, jid, requestData.Text)
		if err != nil {
			fmt.Printf("Error sending message to group %s: %v\n", groupJID, err)
		} else {
			fmt.Printf("Message sent to group %s successfully\n", groupJID)
		}
	}
	*/
	if requestData.Type == "text" {
		err = helpers.SendMessageToPhoneNumber(client, requestData.To, requestData.Text)
		if err != nil {
			// Tangani kesalahan jika gagal mengirim pesan
			helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToSendMessage)
		}
	} else {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidMessageType)
		return
	}
	//fmt.Println("check admin xsilver", adminGroupJIDs)

	response := model.SendMessageResponse{
		ID:     uuid.New().String(),
		From:   requestData.From,
		To:     requestData.To,
		Time:   time.Now().UnixMilli(),
		Status: "delivered",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToMarshalResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func SendMessageBulkHandler(w http.ResponseWriter, r *http.Request) {

	var requestData []model.SendMessageDataRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	var results []map[string]interface{}

	// Send the messages to each recipient
	for _, message := range requestData {
		result := map[string]interface{}{
			"to":      message.To,
			"type":    message.Type,
			"text":    message.Text,
			"caption": message.Caption,
			"url":     message.URL,
			"from":    message.From,
		}

		// Check if any required fields are missing
		if message.Type == "" || message.Text == "" || message.From == "" || message.To == "" {
			result["status"] = "failed"
			results = append(results, result)
			continue
		}

		// Validate 'from' number
		if !helpers.IsValidPhoneNumber(message.From) {
			result["status"] = "failed"
			results = append(results, result)
			continue
		}

		// Validate 'to' number
		if !helpers.IsValidPhoneNumber(message.To) {
			result["status"] = "failed"
			results = append(results, result)
			continue
		}

		// Validate if WhatsApp number is connected and give status failed if number not valid
		if !helpers.IsLoggedInByNumber(client, message.From) {
			result["status"] = "failed"
			results = append(results, result)
			continue
		}

		// Send the message if all checks pass
		if message.Type == "text" {
			err = helpers.SendMessageToPhoneNumber(client, message.To, message.Text)
			if err != nil {
				result["status"] = "failed"
			}
		} else {
			result["status"] = "failed"
		}

		// Add result to results slice
		results = append(results, result)
	}

	// Return the results as the response
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToMarshalResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message, phoneNumber string) {
	errorResponse := response.ErrorResponseNumberPhone{
		StatusCode: statusCode,
		Error:      "Bad Request",
		Message:    "Device with phone: [" + phoneNumber + "] " + message,
	}

	// Convert ErrorResponse to JSON and send it as response
	w.WriteHeader(statusCode)
	jsonBytes, err := json.MarshalIndent(errorResponse, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
