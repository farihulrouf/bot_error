package controllers

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"strings"

	"wagobot.com/errors"
	// "wagobot.com/helpers"
	"wagobot.com/base"
	"wagobot.com/model"
	"wagobot.com/response"
	"go.mau.fi/whatsmeow/types"
)

func GetGroupsHandler(w http.ResponseWriter, r *http.Request) {

	phone := r.URL.Query().Get("phone")
	if phone == "" {
		// helpers.SendErrorResponse(w, http.StatusBadRequest, errors.ErrPhoneNumberRequired)
		base.SetResponse(w, http.StatusBadRequest, errors.ErrPhoneNumberRequired)
		return
	}

	if !base.IsMyNumber(phone) {
		base.SetResponse(w, http.StatusBadRequest, "Missing number")
		return
	}

	var filteredGroups []response.GroupResponse = []response.GroupResponse{}
	mutex.Lock()

	if _, exists := model.Clients[phone]; exists {

		client := model.Clients[phone]

		groups, err := client.Client.GetJoinedGroups()
		if err != nil {
			// helpers.SendErrorResponse(w, http.StatusInternalServerError, errors.ErrFailedToFetchGroups)
			base.SetResponse(w, http.StatusInternalServerError, errors.ErrFailedToFetchGroups)
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

	base.SetResponse(w, http.StatusOK, filteredGroups)
}

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	var params model.JoinGroupRequest

	base.ValidateRequest(r, &params)
	fmt.Println(params)

	if params.Code == "" || params.Phone == "" {
		base.SetResponse(w, http.StatusBadRequest, "phone & code are required")
		return
	}

	phone := params.Phone
	group := params.Code

	if !base.IsMyNumber(phone) {
		base.SetResponse(w, http.StatusBadRequest, "Missing number")
		return
	}

	if _, exists := model.Clients[phone]; exists {
		client := model.Clients[phone].Client

		_, err := client.JoinGroupWithLink(group)
		if err != nil {
			base.SetResponse(w, http.StatusInternalServerError, errors.ErrFailedToJoinGroup)
			return
		}

		base.SetResponse(w, http.StatusOK, "Group joined successfully")
	} else {
		base.SetResponse(w, http.StatusBadRequest, "Invalid account")
	}
}

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	var params model.LeaveGroupRequest

	base.ValidateRequest(r, &params)
	fmt.Println(params)

	if params.Phone == "" || params.GroupID == "" {
		base.SetResponse(w, http.StatusBadRequest, "phone & groupid are required")
		return
	}

	phone := params.Phone

	if !base.IsMyNumber(phone) {
		base.SetResponse(w, http.StatusBadRequest, "Missing number")
		return
	}

	if _, exists := model.Clients[phone]; exists {
		client := model.Clients[phone].Client

		groupJID, err := types.ParseJID(params.GroupID + "@g.us")
		if err != nil {
			base.SetResponse(w, http.StatusBadRequest, errors.ErrInvalidGroupID)
			return
		}

		groupMetadata, err := client.GetGroupInfo(groupJID)
		if err != nil {
			base.SetResponse(w, http.StatusBadRequest, "Invalid group")
			return
		}

		fmt.Println("JID", groupJID)
		fmt.Println("Group", groupMetadata)

		err = client.LeaveGroup(groupJID)
		if err != nil {
			base.SetResponse(w, http.StatusBadRequest, "Leaving group failed")
			return
		}

		base.SetResponse(w, http.StatusOK, "Group joined successfully")
	} else {
		base.SetResponse(w, http.StatusBadRequest, "Invalid account")
	}
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
