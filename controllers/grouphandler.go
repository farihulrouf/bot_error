package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	//"go.mau.fi/whatsmeow"

	"wagobot.com/errors"
	"wagobot.com/helpers"
	"wagobot.com/model"
	"wagobot.com/response"
)

func GetGroupsHandler(w http.ResponseWriter, r *http.Request) {

	phone := r.URL.Query().Get("phone")
	if phone == "" {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrPhoneNumberRequired)
		return
	}

	var filteredGroups []response.GroupResponse
	mutex.Lock()
	for _, client := range clients {
		groups, err := client.GetJoinedGroups()
		if err != nil {
			helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToFetchGroups)
			mutex.Unlock()
			return
		}

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
					Admins:      admins,
					Time:        group.GroupCreated.UnixMilli(),
					Pinned:      false,
					UnreadCount: 30,
				}
				filteredGroups = append(filteredGroups, groupResponse)
			}
		}
	}
	mutex.Unlock()

	response := map[string]interface{}{
		"data": filteredGroups,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToMarshalResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	var value_client = clients["device1"]
	matchFound := false
	var req model.JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//log.Printf("Error decoding request: %v", err)
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidRequestPayload)
		return
	}

	// Ensure all required fields are present
	if req.Code == "" || req.Phone == "" {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "code and phone are required fields")
		return
	}
	//check nunber is login in device

	for key := range clients {
		fmt.Println("Checking key:", key)
		whoami := clients[key].Store.ID.String()
		parts := strings.Split(whoami, ":")
		//fmt.Println("whoami:", whoami)

		if req.Phone == parts[0] {
			fmt.Println("Match found, requestData.From:", req.Phone)
			value_client = clients[key]
			//fmt.Println("whoami:", value_client)
			matchFound = true
			break
		}
	}
	if !matchFound {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "No matching number found")
	}

	// Attempt to join the group with the provided invite link (code)
	groupJID, err := value_client.JoinGroupWithLink(req.Code)
	if err != nil {
		//log.Printf("Error joining group: %v", err)
		helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToJoinGroup)
		return
	}

	// Log success and respond with a success message
	log.Printf("Group joined successfully: %v", groupJID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group joined successfully"})

}

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	/*var req model.LeaveGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//log.Printf("Error decoding request: %v", err)
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidRequestPayload)
		return
	}

	groupJID, err := types.ParseJID(req.GroupID + "@g.us")
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidGroupID)
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
		helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidPhoneNumber)
		return
	}

	log.Printf("Participant %s left group %s successfully, response: %v", req.Phone, req.GroupID, response)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"})
	*/
}

// CreateGroupHandler handles the creation of a new WhatsApp group
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	/*
		var req model.CreateGroupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			//log.Printf("Error decoding request: %v", err)
			helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidRequestPayload)
			//http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if req.Subject == "" {
			helpers.SendErrorResponse(w, http.StatusBadRequest, "Group name is required")
			return
		}

		if len(req.Participants) == 0 {
			helpers.SendErrorResponse(w, http.StatusBadRequest, "At least one participant is required")
			return
		}

		// Convert participant phone numbers to JID
		participants := make([]types.JID, len(req.Participants))
		for i, phone := range req.Participants {
			participantJID, err := types.ParseJID(phone + "@s.whatsapp.net")
			if err != nil {
				helpers.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid phone number: %s", phone))
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
			helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToCreateGroup)
			return
		}

		log.Printf("Group created successfully: %v", groupResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Group created successfully",
			"groupInfo": groupResponse,
		})
	*/
}

// Handler untuk endpoint /api/grouplink
func GetGroupInviteLinkHandler(w http.ResponseWriter, r *http.Request) {
	/*groupID := r.URL.Query().Get("group_id")

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
	*/
}
