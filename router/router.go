package router

import (
	"github.com/gorilla/mux"
	"go.mau.fi/whatsmeow"
	"wagobot.com/controllers"
)

func SetupRouter(client *whatsmeow.Client) *mux.Router {
	r := mux.NewRouter()
	controllers.SetClient(client)

	// Menetapkan penanganan rute untuk endpoint registrasi dan login
	r.HandleFunc("/api/register", controllers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", controllers.LoginHandler).Methods("POST")

	r.HandleFunc("/api/scanqr", controllers.ScanQRHandler).Methods("GET")

	//r.HandleFunc("/api/groups", controllers.CreateGroupHandler).Methods("POST") // New route for creating a group
	r.HandleFunc("/api/groups", controllers.GetGroupsHandler).Methods("GET")
	r.HandleFunc("/api/groups", controllers.JoinGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/messages", controllers.SendMessageGroupHandler).Methods("POST")

	r.HandleFunc("/api/groups/leave", controllers.LeaveGroupHandler).Methods("POST")

	r.HandleFunc("/api/messages", controllers.SendMessageHandler).Methods("POST")
	r.HandleFunc("/api/messages", controllers.RetrieveMessagesHandler).Methods("GET")

	//RetrieveMessagesHandler
	r.HandleFunc("/api/messages/bulk", controllers.SendMessageBulkHandler).Methods("POST")

	r.HandleFunc("/api/results", controllers.GetMessagesHandler).Methods("GET")
	r.HandleFunc("/api/results/{id}", controllers.GetMessagesByIdHandler).Methods("GET")

	r.HandleFunc("/api/logout", controllers.LogoutHandler).Methods("POST") // Add the logout route

	r.HandleFunc("/api/ping", controllers.PingHandler).Methods("GET") // Add the logout route

	r.HandleFunc("/api/system/ver", controllers.VersionHandler).Methods("GET") // Add the logout route

	r.HandleFunc("/api/system/webhook", controllers.SetWebhookHandler).Methods("POST")

	//router.POST("/api/system/webhook", SetWebhookHandler)

	r.HandleFunc("/api/getinfo", controllers.GetInfoHandler).Methods("GET")
	r.HandleFunc("/api/system/devices", controllers.GetDevicesHandler).Methods("GET")

	// Add more routes here if needed
	return r
}
