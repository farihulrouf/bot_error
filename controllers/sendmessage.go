package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"wagobot.com/errors"
	"wagobot.com/helpers"
	"wagobot.com/model"
	"wagobot.com/response"
)

// Pastikan client diimpor dari file yang sesuai
//var client *whatsmeow.Client

func SendMessageGroupHandler(w http.ResponseWriter, r *http.Request) {
	/*
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
	*/
}

// SendMessageHandler handles sending messages.

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var requestData model.SendMessageDataRequest
	var value_client = clients["device1"]
	matchFound := false
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}
	if requestData.To == "" || requestData.Type == "" || requestData.Text == "" || requestData.From == "" {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Missing required fields: 'to', 'type', 'text', or 'from'")
		return
	}

	for key := range clients {
		fmt.Println("Checking key:", key)
		whoami := clients[key].Store.ID.String()
		parts := strings.Split(whoami, ":")
		fmt.Println("whoami:", whoami)
		fmt.Println("whoami:", parts[0])

		if requestData.From == parts[0] {
			fmt.Println("Match found, requestData.From:", requestData.From)
			value_client = clients[key]
			fmt.Println("whoami:", value_client)
			matchFound = true
			break
		}
	}
	if !matchFound {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "No matching number found for requestData.From")
	}

	if requestData.Type == "text" {
		err = helpers.SendMessageToPhoneNumber(value_client, requestData.To, requestData.Text)
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
	var value_client = clients["device1"]
	matchFound := false
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

		for key := range clients {
			fmt.Println("Checking key:", key)
			whoami := clients[key].Store.ID.String()
			parts := strings.Split(whoami, ":")
			fmt.Println("whoami:", whoami)

			if message.From == parts[0] {
				fmt.Println("Match found, requestData.From:", message.From)
				value_client = clients[key]
				fmt.Println("whoami:", value_client)
				matchFound = true
				break
			}
		}
		if !matchFound {
			helpers.SendErrorResponse(w, http.StatusBadRequest, "No matching number found for requestData.From")
		}
		if message.Type == "text" {
			err = helpers.SendMessageToPhoneNumber(value_client, message.To, message.Text)
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
