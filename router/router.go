package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.mau.fi/whatsmeow"
	"wagobot.com/auth"
	"wagobot.com/controllers"
)

func SetupRouter(client *whatsmeow.Client) *mux.Router {
	r := mux.NewRouter()
	controllers.SetClient(client)

	// Menetapkan penanganan rute untuk endpoint registrasi dan login
	r.HandleFunc("/api/register", controllers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", controllers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/scanqr", controllers.ScanQRHandler).Methods("GET")

	// Middleware JWT digunakan untuk semua rute kecuali /api/login dan /api/register /scanqr
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/login" || r.URL.Path == "/api/register" || r.URL.Path == "/api/scanqr" {
				next.ServeHTTP(w, r)
				return
			}
			auth.JWTMiddleware(next).ServeHTTP(w, r)
		})
	})

	//r.HandleFunc("/api/groups", controllers.CreateGroupHandler).Methods("POST") // New route for creating a group
	r.HandleFunc("/api/groups", controllers.GetGroupsHandler).Methods("GET")
	r.HandleFunc("/api/groups", controllers.JoinGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/messages", controllers.SendMessageGroupHandler).Methods("POST")

	r.HandleFunc("/api/groups/leave", controllers.LeaveGroupHandler).Methods("POST")

	r.HandleFunc("/api/messages", controllers.SendMessageHandler).Methods("POST")
	//r.Handle("/api/messages", auth.JWTMiddleware(http.HandlerFunc(controllers.SendMessageHandler))).Methods("POST")
	//r.HandleFunc("/api/messages", auth.JWTMiddleware(controllers.SendMessageHandler)).Methods("POST")
	r.HandleFunc("/api/messages", controllers.RetrieveMessagesHandler).Methods("GET")

	//sendMessage handler
	r.HandleFunc("/api/messages/bulk", controllers.SendMessageBulkHandler).Methods("POST")

	//r.Handle("/api/messages/bulk", auth.JWTMiddleware(http.HandlerFunc(controllers.SendMessageBulkHandler))).Methods("POST")

	r.HandleFunc("/api/result", controllers.GetMessagesHandler).Methods("GET")
	r.HandleFunc("/api/result/{id}", controllers.GetMessagesByIdHandler).Methods("GET")

	r.HandleFunc("/api/logout", controllers.LogoutHandler).Methods("POST") // Add the logout route

	r.HandleFunc("/api/ping", controllers.PingHandler).Methods("GET") // Add the logout route

	r.HandleFunc("/api/system/ver", controllers.VersionHandler).Methods("GET") // Add the logout route

	r.HandleFunc("/api/system/webhook", controllers.SetWebhookHandler).Methods("POST")

	//router.POST("/api/system/webhook", SetWebhookHandler)

	r.HandleFunc("/api/getinfo", controllers.GetInfoHandler).Methods("GET")
	r.HandleFunc("/api/system/devices", controllers.GetDevicesHandler).Methods("GET")

	r.HandleFunc("/api/token", controllers.CreateToken).Methods("POST")

	// Add more routes here if needed
	return r
}
