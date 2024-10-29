package controllers

import (
	// "encoding/json"

	"net/http"
	"strings"

	"wagobot.com/errors"
	// "wagobot.com/helpers"

	"wagobot.com/base"
	"wagobot.com/model"
	"wagobot.com/response"
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
