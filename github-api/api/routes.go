package api

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/repositories", CreateRepository).Methods("POST")
	r.HandleFunc("/repositories/{owner}/{repo}", DeleteRepository).Methods("DELETE")
	r.HandleFunc("/repositories/{user}", ListRepositories).Methods("GET")
	r.HandleFunc("/repositories/{owner}/{repo}/pulls", ListPullRequests).Methods("GET")
}
