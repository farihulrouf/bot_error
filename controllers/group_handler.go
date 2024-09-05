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
	var params model.PhoneCodeRequest

	base.ValidateRequest(r, &params)
	// fmt.Println(params)

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

		groupInfo, err := client.GetGroupInfoFromLink(group)
		if err != nil {
			base.SetResponse(w, http.StatusBadRequest, errors.ErrFailedToJoinGroup)
			return
		}

		_, err = client.JoinGroupWithLink(group)
		if err != nil {
			base.SetResponse(w, http.StatusBadRequest, errors.ErrFailedToJoinGroup)
			return
		}

		groupid := model.GetPhoneNumber(groupInfo.JID.String())

		members := []model.Member{}
		for _, member := range groupInfo.Participants {

			tavatar := saveProfilePicture(client, member.JID)
			tphone  := model.GetPhoneNumber(member.JID.String())
			name := member.DisplayName

			if name == "" {
				contact, err := client.Store.Contacts.GetContact(member.JID)
				if err != nil {
					fmt.Printf("Failed to fetch contact for JID: %s, error: %v\n", member.JID.String(), err)
					continue
				}
				name = contact.PushName
			}

			members = append(members, model.Member{
				ID: tphone,
				Name: name,
				Avatar: tavatar,
				Phone: tphone,
				IsAdmin: member.IsAdmin,
				IsSuperAdmin: member.IsSuperAdmin,
				GroupID: groupid,
			})
		}

		avatar := saveProfilePicture(client, groupInfo.JID)

		payload := model.PayloadWebhook {
			Section: "groups",
			Data: model.Group {
				ID: groupid,
				Name: groupInfo.Name,
				Url: group,
				OwnerID: model.GetPhoneNumber(groupInfo.OwnerJID.String()),
				IsIncognito: groupInfo.IsIncognito,
				IsParent: groupInfo.IsParent,
				Avatar: avatar,
				CreatedTime: groupInfo.GroupCreated.Unix(),
				Members: members,
			},
		}

		err = sendPayloadToWebhook(model.DefaultWebhook, payload)
		if err != nil {
			fmt.Printf("Failed to send payload to webhook: %v\n", err)
		}

		base.SetResponse(w, http.StatusOK, groupInfo)
	} else {
		base.SetResponse(w, http.StatusBadRequest, "Invalid account")
	}
}

func LeaveGroupHandler(w http.ResponseWriter, r *http.Request) {
	var params model.PhoneGroupRequest

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

func MemberGroupHandler(w http.ResponseWriter, r *http.Request) {
	var params model.PhoneGroupRequest

	base.ValidateRequest(r, &params)
	fmt.Println(params)

	if params.GroupID == "" || params.Phone == "" {
		base.SetResponse(w, http.StatusBadRequest, "phone & code are required")
		return
	}

	phone := params.Phone
	groupid := params.GroupID

	if !base.IsMyNumber(phone) {
		base.SetResponse(w, http.StatusBadRequest, "Missing number")
		return
	}

	if _, exists := model.Clients[phone]; exists {
		client := model.Clients[phone].Client

		groupJID, err := types.ParseJID(groupid + "@g.us")
		if err != nil {
			base.SetResponse(w, http.StatusBadRequest, errors.ErrInvalidGroupID)
			return
		}

		groupInfo, err := client.GetGroupInfo(groupJID)
		if err != nil {
			base.SetResponse(w, http.StatusBadRequest, errors.ErrInvalidGroupID)
			return
		}

		members := []model.Member{}
		for _, member := range groupInfo.Participants {

			tavatar := saveProfilePicture(client, member.JID)
			tphone  := model.GetPhoneNumber(member.JID.String())
			name := member.DisplayName

			if name == "" {
				contact, err := client.Store.Contacts.GetContact(member.JID)
				if err != nil {
					fmt.Printf("Failed to fetch contact for JID: %s, error: %v\n", member.JID.String(), err)
					continue
				}
				name = contact.PushName
			}

			members = append(members, model.Member{
				ID: tphone,
				Name: name,
				Avatar: tavatar,
				Phone: tphone,
				IsAdmin: member.IsAdmin,
				IsSuperAdmin: member.IsSuperAdmin,
				GroupID: groupid,
			})
		}

		// fmt.Println(members)

		payload := model.PayloadWebhook {
			Section: "senders",
			Data: members,
		}

		err = sendPayloadToWebhook(model.DefaultWebhook, payload)
		if err != nil {
			fmt.Printf("Failed to send payload to webhook: %v\n", err)
		}

		base.SetResponse(w, http.StatusOK, payload)
	} else {
		base.SetResponse(w, http.StatusBadRequest, "Invalid account")
	}
}

// func ChatGroupHandler(w http.ResponseWriter, r *http.Request) {
// 	var params model.PhoneGroupRequest

// 	base.ValidateRequest(r, &params)
// 	fmt.Println(params)

// 	if params.GroupID == "" || params.Phone == "" {
// 		base.SetResponse(w, http.StatusBadRequest, "phone & code are required")
// 		return
// 	}

// 	phone := params.Phone
// 	groupid := params.GroupID

// 	if !base.IsMyNumber(phone) {
// 		base.SetResponse(w, http.StatusBadRequest, "Missing number")
// 		return
// 	}

// 	if _, exists := model.Clients[phone]; exists {
// 		// client := model.Clients[phone].Client

// 		groupJID, err := types.ParseJID(groupid + "@g.us")
// 		if err != nil {
// 			base.SetResponse(w, http.StatusBadRequest, errors.ErrInvalidGroupID)
// 			return
// 		}

// 		// chat, err := client.Store.GetChat(groupJID)
// 		// if err != nil {
// 		// 	return fmt.Errorf("failed to get chat: %v", err)
// 		// }

// 		fmt.Println(groupJID)

// 		// // Fetch chat history
// 		// ctx := context.Background()
// 		// messages, err := client.GetChatHistory(ctx, groupJID, 100) // Fetch 100 messages
// 		// if err != nil {
// 		// 	log.Fatalf("Failed to get chat history: %v", err)
// 		// }

// 		// // chatJID, err := types.ParseJID(conv.GetId())
// 		// for _, historyMsg := range conv.GetMessages() {
// 		// 	evt, _ := client.ParseWebMessage(groupJID, historyMsg.GetMessage())
// 		// 	// yourNormalEventHandler(evt)
// 		// 	fmt.Println(evt)
// 		// }

// 		// personMsg := map[string][]*events.Message
// 		// evt, err := client.ParseWebMessage(groupJID, historyMsg.GetMessage())
// 		// if err != nil {
// 		// 	// handle
// 		// }
// 		// if !evt.Info.IsFromMe && !evt.Info.IsGroup {// not a group, not sent by me
// 		// 	info, _ := cli.GetUserInfo([]types.JID{evt.Info.Sender})
// 		// 	if contact, ok := contacts[info[evt.Info.Sender]; ok {
// 		// 		msgs, ok := personMsg[contact.PushName]
// 		// 		if !ok {
// 		// 			msgs := []*events.Message{}
// 		// 		}
// 		// 		personMsg[contact.PushName] = append(msgs, evt)
// 		// 	}
// 		// }

// 		// groupInfo, err := client.GetGroupInfo(groupJID)
// 		// if err != nil {
// 		// 	base.SetResponse(w, http.StatusBadRequest, errors.ErrInvalidGroupID)
// 		// 	return
// 		// }

// 		// members := []model.Member{}
// 		// for _, member := range groupInfo.Participants {

// 		// 	tavatar := saveProfilePicture(client, member.JID)
// 		// 	tphone  := model.GetPhoneNumber(member.JID.String())
// 		// 	name := member.DisplayName

// 		// 	if name == "" {
// 		// 		contact, err := client.Store.Contacts.GetContact(member.JID)
// 		// 		if err != nil {
// 		// 			fmt.Printf("Failed to fetch contact for JID: %s, error: %v\n", member.JID.String(), err)
// 		// 			continue
// 		// 		}
// 		// 		name = contact.PushName
// 		// 	}

// 		// 	members = append(members, model.Member{
// 		// 		ID: tphone,
// 		// 		Name: name,
// 		// 		Avatar: tavatar,
// 		// 		Phone: tphone,
// 		// 		IsAdmin: member.IsAdmin,
// 		// 		IsSuperAdmin: member.IsSuperAdmin,
// 		// 		GroupID: groupid,
// 		// 	})
// 		// }

// 		// // fmt.Println(members)

// 		// payload := model.PayloadWebhook {
// 		// 	Section: "senders",
// 		// 	Data: members,
// 		// }

// 		// err = sendPayloadToWebhook(model.DefaultWebhook, payload)
// 		// if err != nil {
// 		// 	fmt.Printf("Failed to send payload to webhook: %v\n", err)
// 		// }

// 		base.SetResponse(w, http.StatusOK, messages)
// 	} else {
// 		base.SetResponse(w, http.StatusBadRequest, "Invalid account")
// 	}
// }

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
