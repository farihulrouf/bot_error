package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/groups", createGroupHandler).Methods("POST")
	r.HandleFunc("/api/groups", getGroupsHandler).Methods("GET")
	return r
}
