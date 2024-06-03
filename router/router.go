package router

import (
	"github.com/gorilla/mux"
	"go.mau.fi/whatsmeow"
	"wagobot.com/controllers"
)

func SetupRouter(client *whatsmeow.Client) *mux.Router {
	r := mux.NewRouter()
	controllers.SetClient(client)
	r.HandleFunc("/api/groups", controllers.GetGroupsHandler).Methods("GET")
	r.HandleFunc("/api/groups", controllers.CreateGroupHandler).Methods("POST") // New route for creating a group
	r.HandleFunc("/api/groups/messages", controllers.SendMessageGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/join", controllers.JoinGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/leave", controllers.LeaveGroupHandler).Methods("POST")

	r.HandleFunc("/api/messages", controllers.SendMessageHandler).Methods("POST")
	r.HandleFunc("/api/messages/bulk", controllers.SendMessageBulkHandler).Methods("POST")

	r.HandleFunc("/api/results", controllers.GetMessagesHandler).Methods("GET")
	r.HandleFunc("/api/results/{id}", controllers.GetMessagesByIdHandler).Methods("GET")

	r.HandleFunc("/api/logout", controllers.LogoutHandler).Methods("POST") // Add the logout route

	// Add more routes here if needed
	return r
}
