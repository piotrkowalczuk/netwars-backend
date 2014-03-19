package user

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/binding"
)

func CreateRoute(router martini.Router) () {
	router.Post("/user",
		binding.Json(APICredentials{}),
		AuthenticationMiddleware,
		createHandler,
	)
	router.Get("/user/:id",
		binding.Json(APICredentials{}),
		AuthenticationMiddleware,
		readHandler,
	)
	router.Put("/user/:id",
		binding.Json(APICredentials{}),
		AuthenticationMiddleware,
		updateHandler,
	)
	router.Delete("/user/:id",
		binding.Json(APICredentials{}),
		AuthenticationMiddleware,
		deleteHandler,
	)

	router.Post("/login", binding.Json(LoginCredentials{}), loginHandler)
	router.Post("/logout", binding.Json(APICredentials{}), logoutHandler)
	router.Post("/register", binding.Json(User{}), logoutHandler)
}
