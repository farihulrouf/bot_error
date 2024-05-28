package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"wagobot.com/model"
)

var client *whatsmeow.Client // Define client variable
func Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group created successfully"})
}
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.CreateGroupRequest // Update to use the model package
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	participantJID, err := types.ParseJID(req.Phone + "@s.whatsapp.net")
	if err != nil {
		log.Printf("Error parsing JID: %v", err)
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	reqCreateGroup := whatsmeow.ReqCreateGroup{
		//Participants: []whatsmeow.JID{participantJID},
		Participants: []types.JID{participantJID},
	}

	groupResponse, err := client.CreateGroup(reqCreateGroup)
	if err != nil {
		log.Printf("Error creating group: %v", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	log.Printf("Group created successfully: %v", groupResponse)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group created successfully"})
}

func GetGroupsHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	groups, err := client.GetJoinedGroups()
	if err != nil {
		http.Error(w, "Failed to get groups", http.StatusInternalServerError)
		return
	}

	groupList := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		fmt.Printf("Group ID: %s, Name: %s\n", group.JID.String(), group.Name)
		groupInfo := map[string]interface{}{
			"JID":  group.JID.String(),
			"Name": group.Name,
		}
		groupList = append(groupList, groupInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupList)
}

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LeaveGroupRequest // Update to use the model package
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	groupJID, err := types.ParseJID(req.GroupID + "@g.us")
	if err != nil {
		log.Printf("Error parsing Group JID: %v", err)
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	participantJID, err := types.ParseJID(req.Phone + "@s.whatsapp.net")
	if err != nil {
		log.Printf("Error parsing Participant JID: %v", err)
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	//  Dengan asumsi metode untuk memperbarui peserta grup adalah UpdateGroupParticipants
	response, err := client.UpdateGroupParticipants(groupJID, []types.JID{participantJID}, "remove")
	//response, err := client.UpdateGroupParticipants(groupJID, []whatsmeow.JID{participantJID}, "remove")
	if err != nil {
		log.Printf("Error leaving group: %v", err)
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	log.Printf("Participant %s left group %s successfully, response: %v", req.Phone, req.GroupID, response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"})
}
