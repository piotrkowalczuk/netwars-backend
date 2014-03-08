package server

import (
	"github.com/piotrkowalczuk/netwars-backend/user"
	"github.com/gorilla/mux"
)

func CreateRoute() (*mux.Router) {
	router := mux.NewRouter()
	user.CreateRoute(router)

	return router
}
