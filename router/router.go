package router

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	// sesuaikan dengan path yang sesuai
	//"go.mau.fi/whatsmeow"

	"wagobot.com/base"
	"wagobot.com/controllers"
)

func SetupRouter() *mux.Router {
	
	r := mux.NewRouter()

	// Middleware JWT digunakan untuk semua rute kecuali /api/login dan /api/register /scanqr
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/login" || r.URL.Path == "/api/register" ||
				strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}
			base.JWTMiddleware(next).ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/", controllers.VersionHandler).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler) // not ok, file not found

	r.HandleFunc("/api/login", controllers.LoginHandler).Methods("POST") // ok
	r.HandleFunc("/api/logout", controllers.LogoutHandler).Methods("POST") // not ok, response kosong tidak berfungsi
	
	r.HandleFunc("/api/ping", controllers.PingHandler).Methods("GET") // response ganti ke json {status: "online/offline"}
	r.HandleFunc("/api/register", controllers.RegisterHandler).Methods("POST") // not ok, tidak ada validasi parameter user
	r.HandleFunc("/api/token", controllers.CreateToken).Methods("POST") // not ok, token tidak bisa digunakan

	r.HandleFunc("/api/system/logout/{phone}", controllers.RemoveClient).Methods("DELETE")
	r.HandleFunc("/api/system/devices", controllers.CreateDevice).Methods("GET") // not ok, hang
	r.HandleFunc("/api/system/ver", controllers.VersionHandler).Methods("GET") // ok
	r.HandleFunc("/api/system/webhook", controllers.WebhookHandler).Methods("POST")
	// router.POST("/api/system/webhook", SetWebhookHandler)

	r.HandleFunc("/api/webhook/update", controllers.UpdateWbhookURLHandler).Methods("PUT")

	r.HandleFunc("/api/user", controllers.GetUserHandler).Methods("GET") // ok
	r.HandleFunc("/api/user", controllers.UserUpdateHandler).Methods("PUT")
	r.HandleFunc("/api/user/login", controllers.LoginHandler).Methods("POST") // ok
	// r.HandleFunc("/api/user/detail", controllers.GetUserHandler).Methods("GET") // ok
	// r.HandleFunc("/api/user/update", controllers.UserUpdateHandler).Methods("PUT")
	
	r.HandleFunc("/api/group/invite", controllers.GetGroupInviteLinkHandler).Methods("GET")

	r.HandleFunc("/api/groups", controllers.GetGroupsHandler).Methods("GET")
	r.HandleFunc("/api/group/join", controllers.JoinGroupHandler).Methods("POST")
	r.HandleFunc("/api/group/leave", controllers.LeaveGroupHandler).Methods("POST")
	// r.HandleFunc("/api/groups", controllers.JoinGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups/messages", controllers.SendMessageGroupHandler).Methods("POST")
	// r.HandleFunc("/api/groups/leave", controllers.LeaveGroupHandler).Methods("POST")

	r.HandleFunc("/api/messages", controllers.GetSearchMessagesHandler).Methods("GET")
	r.HandleFunc("/api/messages", controllers.SendMessageHandler).Methods("POST")
	r.HandleFunc("/api/messages/bulk", controllers.SendMessageBulkHandler).Methods("POST")
	
	r.HandleFunc("/api/result", controllers.GetMessagesHandler).Methods("GET")
	r.HandleFunc("/api/result/{id}", controllers.GetMessagesByIdHandler).Methods("GET")
	
	//r.HandleFunc("/api/get/client", controllers.GetClientByDeviceNameHandler).Methods("GET")
	//r.HandleFunc("/status/qr/list", controllers.GetConnectedClientsList).Methods("GET")
	//r.HandleFunc("/api/token", controllers.CreateToken).Methods("POST")
	//r.HandleFunc("/api/messages/images", controllers.SendImageHandler).Methods("POST")

	r.HandleFunc("/api/devices", controllers.GetDevicesHandler).Methods("GET")
	r.HandleFunc("/api/device/scan", controllers.ScanDeviceHandler).Methods("GET")
	// r.HandleFunc("/api/device/{id}", controllers.GetDevicesHandler).Methods("GET")
	// r.HandleFunc("/api/device/{id}", controllers.GetDevicesHandler).Methods("PUT")
	// r.HandleFunc("/api/device/{id}", controllers.GetDevicesHandler).Methods("DELETE")


	return r
}
