package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"wagobot.com/controllers"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/test", controllers.Test).Methods("GET")
	r.HandleFunc("/api/groups", controllers.CreateGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups", controllers.GetGroupsHandler).Methods("GET")
	r.HandleFunc("/api/groups/leave", controllers.LeaveGroupHandler).Methods("POST")

	return r
}

func RunServer() {
	r := NewRouter()
	http.Handle("/", r)

	// Start server
	http.ListenAndServe(":8080", nil)
}
