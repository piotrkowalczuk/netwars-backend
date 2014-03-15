package user

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/binding"
)

func CreateRoute(router martini.Router) () {
	router.Post("/user", AuthenticationMiddleware, create)
	router.Get("/user/:id", AuthenticationMiddleware, read)
	router.Put("/user/:id", AuthenticationMiddleware, update)
	router.Delete("/user/:id", AuthenticationMiddleware, delete)

	router.Post("/login", binding.Json(LoginCredentials{}), login)
	router.Post("/logout", binding.Json(APICredentials{}), logout)
}
