package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func CreateUserRoute(router martini.Router) () {
	router.Get("/user/:id", getUserHandler)
	router.Get("/users/online", getOnlineUsersHandler)
	router.Post("/login", binding.Json(LoginCredentials{}), loginHandler)
	router.Post(
		"/logout",
		binding.Json(APICredentials{}),
		AuthenticationMiddleware,
		logoutHandler,
	)
	router.Post("/register", binding.Json(UserRegistration{}), registerHandler)
}
