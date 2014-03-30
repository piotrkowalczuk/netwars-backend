package user

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/binding"
)

func CreateRoute(router martini.Router) () {
	router.Get("/user/:id", getUserHandler)
	router.Post("/login", binding.Json(LoginCredentials{}), loginHandler)
	router.Post("/logout", binding.Json(APICredentials{}), logoutHandler)
	router.Post("/register", binding.Json(User{}), registerHandler)
}
