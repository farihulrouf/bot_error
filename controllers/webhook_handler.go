package controllers

import (
	"errors"
	"fmt"
	// "strings"
	"net/http"
	"database/sql"
	// "encoding/json"
	// "golang.org/x/crypto/bcrypt"
	"wagobot.com/base"
	"wagobot.com/db"
	"wagobot.com/model"
	// "wagobot.com/response"
)

func GetWebhookHandler(w http.ResponseWriter, r *http.Request) {
	username := base.CurrentUser.Username
	user, err := db.GetUserByID(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.SetResponse(w, http.StatusNotFound, "User not found")
		} else {
			base.SetResponse(w, http.StatusInternalServerError, "Failed to get user data")
		}
		return
	}

	base.SetResponse(w, http.StatusOK, user.Url)
}

func UpdateWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var params model.ChangeWebhookRequest

	base.ValidateRequest(r, &params)
	fmt.Println(params)

	if params.Url == "" {
		base.SetResponse(w, http.StatusBadRequest, "url parameter is required")
		return
	}

	username := base.CurrentUser.Username

	err := db.UpdateUserURLWebhook(username, params.Url)

	if err != nil {
		base.SetResponse(w, http.StatusInternalServerError, "update error")
		return
	}

	// // update active clients
	// for key, client := range model.Clients {
	// 	if strings.HasPrefix(key, "DEV") {
	// 		continue
	// 	}
	// 	if base.CurrentUser.ID == client.User {
	// 		// ermove old client
	// 		delete(model.Clients, key)
	// 		// add new client
	// 		clientLog := waLog.Stdout("Client", "DEBUG", true)
	// 		client := whatsmeow.NewClient(device, clientLog)
	// 		controllers.AddClient(base.CurrentUser.ID, params.Url, key, client, 0)
	// 	}
	// }

	base.SetResponse(w, http.StatusOK, "Webhook changed successfully")
}

