package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	//"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"wagobot.com/errors"
	"wagobot.com/helpers"
	"wagobot.com/model"
	"wagobot.com/response"
)

func GetGroupsHandler(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	groups, err := client.GetJoinedGroups()
	if err != nil {
		http.Error(w, "Failed to fetch joined groups", http.StatusInternalServerError)
		return
	}
	var filteredGroups []response.GroupResponse
	for _, group := range groups {
		var isMember bool
		var members []string
		var admins []string

		for _, member := range group.Participants {
			if member.JID.User == phone {
				isMember = true
			}
			members = append(members, member.JID.User)

			// Check if the member is an admin
			if member.IsAdmin {
				admins = append(admins, member.JID.User)
			}
		}

		if isMember {
			groupID := strings.TrimSuffix(group.JID.String(), "@g.us")
			groupResponse := response.GroupResponse{
				ID:          groupID,
				Type:        "group",
				Description: group.Name,
				Members:     members,
				Admins:      admins, //time.Now().UnixMilli()
				Time:        group.GroupCreated.UnixMilli(),
				Pinned:      false,
				UnreadCount: 30,
			}
			filteredGroups = append(filteredGroups, groupResponse)
		}
	}

	response := map[string]interface{}{
		"data": filteredGroups,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToMarshalResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure all required fields are present
	if req.Code == "" || req.Phone == "" {
		http.Error(w, "code and phone are required fields", http.StatusBadRequest)
		return
	}

	if !helpers.IsLoggedInByNumber(client, req.Phone) {
		sendErrorResponse(w, http.StatusBadRequest, "not ready or not available. Please pairing the device", req.Phone)
		return
	}

	// Attempt to join the group with the provided invite link (code)
	groupJID, err := client.JoinGroupWithLink(req.Code)
	if err != nil {
		log.Printf("Error joining group: %v", err)
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	// Log success and respond with a success message
	log.Printf("Group joined successfully: %v", groupJID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group joined successfully"})
}

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LeaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	/*
		if !helpers.IsLoggedInByNumber(client, requestData.From) {
			http.Error(w, "Sender number is not connected to WhatsApp", http.StatusBadRequest)
			return
		}
	*/

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
	if err != nil {
		log.Printf("Error leaving group: %v", err)
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	log.Printf("Participant %s left group %s successfully, response: %v", req.Phone, req.GroupID, response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"})
}

// CreateGroupHandler handles the creation of a new WhatsApp group
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var req model.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Subject == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	if len(req.Participants) == 0 {
		http.Error(w, "At least one participant is required", http.StatusBadRequest)
		return
	}

	// Convert participant phone numbers to JID
	participants := make([]types.JID, len(req.Participants))
	for i, phone := range req.Participants {
		participantJID, err := types.ParseJID(phone + "@s.whatsapp.net")
		if err != nil {
			log.Printf("Error parsing JID for phone %s: %v", phone, err)
			http.Error(w, fmt.Sprintf("Invalid phone number: %s", phone), http.StatusBadRequest)
			return
		}
		participants[i] = participantJID
	}

	// Hypothetical field for group creation; replace with actual field names from the library
	reqCreateGroup := whatsmeow.ReqCreateGroup{
		// Try possible field names; here we assume 'Name' and 'Participants'
		Name:         req.Subject, // Hypothetical field
		Participants: participants,
	}

	groupResponse, err := client.CreateGroup(reqCreateGroup)
	if err != nil {
		log.Printf("Error creating group: %v", err)
		http.Error(w, fmt.Sprintf("Failed to create group: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Group created successfully: %v", groupResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Group created successfully",
		"groupInfo": groupResponse,
	})
}

// Handler untuk endpoint /api/grouplink
func GetGroupInviteLinkHandler(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")

	// Mendapatkan nilai reset dari parameter URL (atau dari body, sesuai kebutuhan)
	reset := r.URL.Query().Get("reset")

	// Lakukan validasi parameter jika diperlukan
	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Konversi reset menjadi boolean
	resetBool := false
	if reset == "true" {
		resetBool = true
	}
	groupJID, err := types.ParseJID(groupID + "@g.us")

	inviteLink, err := client.GetGroupInviteLink(groupJID, resetBool)
	if err != nil {
		http.Error(w, "Failed to get group invite link", http.StatusInternalServerError)
		return
	}

	// Mengirimkan tautan undangan sebagai respons
	response := map[string]string{"group_invite_link": inviteLink}
	json.NewEncoder(w).Encode(response)
}
