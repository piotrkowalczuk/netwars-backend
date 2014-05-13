package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func CreateUserRoute(router martini.Router) () {
	// GET
	router.Get("/user/stream",
		binding.Form(APICredentials{}),
		getUserStreamHandler,
	)
	router.Get("/user/:id", getUserHandler)
	router.Get("/users/online", getOnlineUsersHandler)
	// POST
	router.Post("/login", binding.Json(LoginCredentials{}), loginHandler)
	router.Post(
		"/logout",
		binding.Json(APICredentials{}),
		AuthenticationMiddleware,
		logoutHandler,
	)
	router.Post("/register", binding.Json(UserRegistration{}), registerHandler)
	// PUT
	router.Put(
		"/users/stream",
		binding.Json(StreamRequest{}),
		binding.Form(APICredentials{}),
		putUserStreamHandler,
	)
}
