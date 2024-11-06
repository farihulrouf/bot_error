package router

import (

	// "encoding/json"

	"github.com/gorilla/mux"

	// sesuaikan dengan path yang sesuai
	//"go.mau.fi/whatsmeow"

	"wagobot.com/controllers"
)

func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	// Middleware JWT digunakan untuk semua rute kecuali /api/login dan /api/register /scanqr
	/*r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			fmt.Printf("%s %s\n", r.Method, r.URL.Path)
			fmt.Printf("Params: %v\n", r.URL.Query())
			fmt.Printf("Forms: %v\n", r.Form)

			if r.URL.Path == "/api/login" || r.URL.Path == "/api/register" ||
				strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}
			base.JWTMiddleware(next).ServeHTTP(w, r)
		})
	})*/

	//r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler) // not ok, file not found

	r.HandleFunc("/api/login", controllers.LoginHandler).Methods("POST")   // ok
	r.HandleFunc("/api/logout", controllers.LogoutHandler).Methods("POST") // not ok, response kosong tidak berfungsi
	// response ganti ke json {status: "online/offline"}
	r.HandleFunc("/api/register", controllers.RegisterHandler).Methods("POST") // not ok, tidak ada validasi parameter user

	r.HandleFunc("/api/system/logout/{phone}", controllers.RemoveClient).Methods("DELETE")

	r.HandleFunc("/api/bot-report", controllers.SendMessageGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups", controllers.GetGroupsHandler).Methods("GET")

	// r.HandleFunc("/api/groups/leave", controllers.LeaveGroupHandler).Methods("POST")

	r.HandleFunc("/api/device/scan", controllers.ScanDeviceHandler).Methods("GET")

	r.HandleFunc("/api/device/remove", controllers.RemoveDeviceHandler).Methods("DELETE")

	// r.HandleFunc("/api/group/messages", controllers.MemberGroupHandler).Methods("POST")

	return r
}
