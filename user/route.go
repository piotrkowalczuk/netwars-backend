package user

import "github.com/gorilla/mux"

func CreateRoute(parentRoute *mux.Router) (childRoute *mux.Router) {
	childRoute = parentRoute.PathPrefix("/user").Subrouter()

	childRoute.HandleFunc("/", create).Methods("POST")
	childRoute.HandleFunc("/{id:[0-9]+}", read).Methods("GET")
	childRoute.HandleFunc("/{id:[0-9]+}", update).Methods("PUT")
	childRoute.HandleFunc("/{id:[0-9]+}", delete).Methods("DELETE")

	return
}
