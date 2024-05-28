package main

import (
	"encoding/json"
	"net/http"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func createGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	participantJID, err := types.ParseJID(req.Phone + "@s.whatsapp.net")
	if err != nil {
		http.Error(w, "Invalid phone number", http.StatusBadRequest)
		return
	}

	reqCreateGroup := whatsmeow.ReqCreateGroup{
		Participants: []types.JID{participantJID},
	}

	_, err = client.CreateGroup(reqCreateGroup)
	if err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group created successfully"})
}

func getGroupsHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := client.GetJoinedGroups()
	if err != nil {
		http.Error(w, "Failed to get groups", http.StatusInternalServerError)
		return
	}

	groupList := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		groupInfo := map[string]interface{}{
			"JID":  group.JID.String(),
			"Name": group.Name,
		}
		groupList = append(groupList, groupInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupList)
}
