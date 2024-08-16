package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"wagobot.com/errors"
	"wagobot.com/helpers"
	"wagobot.com/model"
)

// Pastikan client diimpor dari file yang sesuai
//var client *whatsmeow.Client

func SendMessageGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.SendMessageDataRequest
	var value_client = clients["device1"]
	matchFound := false
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
		http.Error(w, fmt.Sprintf(errors.ErrInvalidRecipient, err), http.StatusBadRequest)
		return
	}

	for key := range clients {
		//fmt.Println("Checking key:", key)
		whoami := clients[key].Store.ID.String()
		parts := strings.Split(whoami, ":")
		//fmt.Println("whoami:", parts[0])

		if req.From == parts[0] {
			//fmt.Println("Match found, requestData.From:", req.From)
			value_client = clients[key]
			matchFound = true
			break
		}
	}

	if !matchFound {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "No matching number found for requestData.From")
		return
	}

	// Send the message
	if err := helpers.SendMessage(value_client, jid, req); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf(errors.ErrInvalidMessageType, err))
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	//fmt.Fprintf(w, "Message sent to: %s", req.To)

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
		//fmt.Println("Checking key:", key)
		whoami := clients[key].Store.ID.String()
		parts := strings.Split(whoami, ":")
		//fmt.Println("whoami:", parts[0])

		if requestData.From == parts[0] {
			//fmt.Println("Match found, requestData.From:", requestData.From)
			value_client = clients[key]
			matchFound = true
			break
		}
	}

	if !matchFound {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "No matching number found for requestData.From")
		return
	}

	if requestData.Type == "text" {
		err = helpers.SendMessageToPhoneNumber(value_client, requestData.To, requestData.Text)
		if err != nil {
			helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToSendMessage)
			return
		}
	} else if requestData.Type == "image" {
		var imageBytes []byte
		var err error // Deklarasikan variabel err di sini

		if helpers.IsValidURL(requestData.URL) {
			fileType, detectErr := helpers.DetectFileTypeByContentType(requestData.URL)
			fmt.Println("cek data", fileType)
			if detectErr != nil {
				fmt.Println("Error detecting file type:", detectErr)
				http.Error(w, "Error detecting file type", http.StatusInternalServerError)
				return
			}

			imageBytes, err = helpers.DownloadFile(requestData.URL)
			if err != nil {
				fmt.Println("Error downloading file:", err)
				http.Error(w, "Error downloading file", http.StatusInternalServerError)
				return
			}
		} else if helpers.IsBase64(requestData.URL) {
			// Jangan deklarasikan variabel baru, cukup gunakan yang sudah ada
			var decodeErr error
			imageBytes, decodeErr = base64.StdEncoding.DecodeString(requestData.URL)
			if decodeErr != nil {
				fmt.Println("Error decoding Base64 string:", decodeErr)
				http.Error(w, "Error decoding Base64 string", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Invalid URL and not a Base64 or url string", http.StatusBadRequest)
			return
		}

		imageMsg, err := helpers.UploadImageAndCreateMessage(value_client, imageBytes, requestData.Caption, requestData.Type)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Kirim pesan gambar
		err = helpers.SendImageToPhoneNumber(value_client, requestData.To, imageMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else if requestData.Type == "doc" {
		var docBytes []byte
		var err error // Deklarasikan variabel err di sini

		if helpers.IsValidURL(requestData.URL) {
			fileType, detectErr := helpers.DetectFileTypeByContentType(requestData.URL)
			fmt.Println("cek data", fileType)
			if detectErr != nil {
				fmt.Println("Error detecting file type:", detectErr)
				http.Error(w, "Error detecting file type", http.StatusInternalServerError)
				return
			}

			docBytes, err = helpers.DownloadFile(requestData.URL)
			if err != nil {
				fmt.Println("Error downloading file:", err)
				http.Error(w, "Error downloading file", http.StatusInternalServerError)
				return
			}
		} else if helpers.IsBase64(requestData.URL) {
			// Jangan deklarasikan variabel baru, cukup gunakan yang sudah ada
			var decodeErr error
			docBytes, decodeErr = base64.StdEncoding.DecodeString(requestData.URL)
			if decodeErr != nil {
				fmt.Println("Error decoding Base64 string:", decodeErr)
				http.Error(w, "Error decoding Base64 string", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Invalid URL and not a Base64 or url string", http.StatusBadRequest)
			return
		}

		docMsg, err := helpers.UploadDocAndCreateMessage(value_client, docBytes, requestData.Caption, requestData.Type)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Kirim pesan doc
		err = helpers.SendDocToPhoneNumber(value_client, requestData.To, docMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if requestData.Type == "video" {
		var vidBytes []byte
		var err error // Deklarasikan variabel err di sini

		if helpers.IsValidURL(requestData.URL) {
			fileType, detectErr := helpers.DetectFileTypeByContentType(requestData.URL)
			fmt.Println("cek data", fileType)
			if detectErr != nil {
				fmt.Println("Error detecting file type:", detectErr)
				http.Error(w, "Error detecting file type", http.StatusInternalServerError)
				return
			}

			vidBytes, err = helpers.DownloadFile(requestData.URL)
			if err != nil {
				fmt.Println("Error downloading file:", err)
				http.Error(w, "Error downloading file", http.StatusInternalServerError)
				return
			}
		} else if helpers.IsBase64(requestData.URL) {
			// Jangan deklarasikan variabel baru, cukup gunakan yang sudah ada
			var decodeErr error
			vidBytes, decodeErr = base64.StdEncoding.DecodeString(requestData.URL)
			if decodeErr != nil {
				fmt.Println("Error decoding Base64 string:", decodeErr)
				http.Error(w, "Error decoding Base64 string", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Invalid URL and not a Base64 or url string", http.StatusBadRequest)
			return
		}

		vidMsg, err := helpers.UploadVideoAndCreateMessage(value_client, vidBytes, requestData.Caption, requestData.Type)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Kirim pesan video
		err = helpers.SendVideoToPhoneNumber(value_client, requestData.To, vidMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidMessageType)
		return
	}

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
	allSucceeded := true
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
			//fmt.Println("Checking key:", key)
			whoami := clients[key].Store.ID.String()
			parts := strings.Split(whoami, ":")
			//fmt.Println("whoami:", whoami)

			if message.From == parts[0] {
				fmt.Println("Match found, requestData.From:", message.From)
				value_client = clients[key]
				//fmt.Println("whoami:", value_client)
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
				allSucceeded = false // Set flag ke false jika tipe bukan "text"
			}

		} else if message.Type == "image" {
			var imageBytes []byte
			var err error // Deklarasikan variabel err di sini

			if helpers.IsValidURL(message.URL) {
				fileType, detectErr := helpers.DetectFileTypeByContentType(message.URL)
				fmt.Println("cek data", fileType)
				if detectErr != nil {
					fmt.Println("Error detecting file type:", detectErr)
					http.Error(w, "Error detecting file type", http.StatusInternalServerError)
					return
				}

				imageBytes, err = helpers.DownloadFile(message.URL)
				if err != nil {
					fmt.Println("Error downloading file:", err)
					http.Error(w, "Error downloading file", http.StatusInternalServerError)
					return
				}
			} else if helpers.IsBase64(message.URL) {
				// Jangan deklarasikan variabel baru, cukup gunakan yang sudah ada
				var decodeErr error
				imageBytes, decodeErr = base64.StdEncoding.DecodeString(message.URL)
				if decodeErr != nil {
					fmt.Println("Error decoding Base64 string:", decodeErr)
					http.Error(w, "Error decoding Base64 string", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid URL and not a Base64 or url string", http.StatusBadRequest)
				return
			}

			imageMsg, err := helpers.UploadImageAndCreateMessage(value_client, imageBytes, message.Text, message.Type)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Kirim pesan gambar
			err = helpers.SendImageToPhoneNumber(value_client, message.To, imageMsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if message.Type == "doc" {
			var docBytes []byte
			var err error // Deklarasikan variabel err di sini

			if helpers.IsValidURL(message.URL) {
				fileType, detectErr := helpers.DetectFileTypeByContentType(message.URL)
				fmt.Println("cek data", fileType)
				if detectErr != nil {
					fmt.Println("Error detecting file type:", detectErr)
					http.Error(w, "Error detecting file type", http.StatusInternalServerError)
					return
				}

				docBytes, err = helpers.DownloadFile(message.URL)
				if err != nil {
					fmt.Println("Error downloading file:", err)
					http.Error(w, "Error downloading file", http.StatusInternalServerError)
					return
				}
			} else if helpers.IsBase64(message.URL) {
				// Jangan deklarasikan variabel baru, cukup gunakan yang sudah ada
				var decodeErr error
				docBytes, decodeErr = base64.StdEncoding.DecodeString(message.URL)
				if decodeErr != nil {
					fmt.Println("Error decoding Base64 string:", decodeErr)
					http.Error(w, "Error decoding Base64 string", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid URL and not a Base64 or url string", http.StatusBadRequest)
				return
			}

			docMsg, err := helpers.UploadDocAndCreateMessage(value_client, docBytes, message.Caption, message.Type)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Kirim pesan doc
			err = helpers.SendDocToPhoneNumber(value_client, message.To, docMsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if message.Type == "video" {
			var vidBytes []byte
			var err error // Deklarasikan variabel err di sini

			if helpers.IsValidURL(message.URL) {
				fileType, detectErr := helpers.DetectFileTypeByContentType(message.URL)
				fmt.Println("cek data", fileType)
				if detectErr != nil {
					fmt.Println("Error detecting file type:", detectErr)
					http.Error(w, "Error detecting file type", http.StatusInternalServerError)
					return
				}

				vidBytes, err = helpers.DownloadFile(message.URL)
				if err != nil {
					fmt.Println("Error downloading file:", err)
					http.Error(w, "Error downloading file", http.StatusInternalServerError)
					return
				}
			} else if helpers.IsBase64(message.URL) {
				// Jangan deklarasikan variabel baru, cukup gunakan yang sudah ada
				var decodeErr error
				vidBytes, decodeErr = base64.StdEncoding.DecodeString(message.URL)
				if decodeErr != nil {
					fmt.Println("Error decoding Base64 string:", decodeErr)
					http.Error(w, "Error decoding Base64 string", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Invalid URL and not a Base64 or url string", http.StatusBadRequest)
				return
			}

			vidMsg, err := helpers.UploadVideoAndCreateMessage(value_client, vidBytes, message.Caption, message.Type)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Kirim pesan doc
			err = helpers.SendVideoToPhoneNumber(value_client, message.To, vidMsg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			result["status"] = "failed"
			allSucceeded = false // Set flag ke false jika tipe bukan "text"
		}

		// Add result to results slice
		results = append(results, result)
	}
	if allSucceeded {
		response := map[string]bool{"queue": true}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToMarshalResponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		jsonResponse, err := json.Marshal(results)
		if err != nil {
			helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToMarshalResponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}
