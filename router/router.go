package router

import (
	"net/http"

	"github.com/gorilla/mux"
	//"go.mau.fi/whatsmeow"

	"wagobot.com/auth"
	"wagobot.com/controllers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	//controllers.SetClient(client)

	// Menetapkan penanganan rute untuk endpoint registrasi dan login
	r.HandleFunc("/api/register", controllers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", controllers.LoginHandler).Methods("POST")
	//r.HandleFunc("/api/scanqr/{device}", controllers.ScanQRHandler).Methods("GET")

	// Middleware JWT digunakan untuk semua rute kecuali /api/login dan /api/register /scanqr
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/login" || r.URL.Path == "/api/register" {
				next.ServeHTTP(w, r)
				return
			}
			auth.JWTMiddleware(next).ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/api/groups", controllers.GetGroupsHandler).Methods("GET")
	r.HandleFunc("/api/groups", controllers.JoinGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/messages", controllers.SendMessageGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/leave", controllers.LeaveGroupHandler).Methods("POST")
	r.HandleFunc("/api/messages", controllers.SendMessageHandler).Methods("POST")

	r.HandleFunc("/api/messages/bulk", controllers.SendMessageBulkHandler).Methods("POST")
	r.HandleFunc("/api/messages", controllers.GetSearchMessagesHandler).Methods("GET")
	r.HandleFunc("/api/token", controllers.CreateToken).Methods("POST")
	r.HandleFunc("/api/result", controllers.GetMessagesHandler).Methods("GET")
	r.HandleFunc("/api/result/{id}", controllers.GetMessagesByIdHandler).Methods("GET")

	r.HandleFunc("/api/logout", controllers.LogoutHandler).Methods("POST") // Add the logout route

	r.HandleFunc("/api/ping", controllers.PingHandler).Methods("GET") // Add the logout route

	r.HandleFunc("/api/system/ver", controllers.VersionHandler).Methods("GET") // Add the logout route

	r.HandleFunc("/api/system/webhook", controllers.WebhookHandler).Methods("POST")

	r.HandleFunc("/api/group/invite", controllers.GetGroupInviteLinkHandler).Methods("GET")

	//router.POST("/api/system/webhook", SetWebhookHandler)

	//r.HandleFunc("/api/getinfo", controllers.GetInfoHandler).Methods("GET")
	r.HandleFunc("/api/system/devices", controllers.CreateDevice).Methods("GET")
	r.HandleFunc("/api/get/client", controllers.GetClientByDeviceNameHandler).Methods("GET")

	//r.HandleFunc("/api/system", controllers.ListDevices).Methods("GET")
	//r.HandleFunc("/status/qr/list", controllers.GetConnectedClientsList).Methods("GET")

	///r.HandleFunc("/api/token", controllers.CreateToken).Methods("POST")

	r.HandleFunc("/api/webhook/update", controllers.UpdateUserURLHandler).Methods("PUT")
	r.HandleFunc("/api/user/detail", controllers.GetUserHandler).Methods("GET")
	//GetUserHandler
	// Add more routes here if needed
	return r
}
