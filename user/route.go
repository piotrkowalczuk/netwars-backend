package user

import "github.com/codegangsta/martini"

func CreateRoute(m *martini.Martini) () {
	childRoute := martini.NewRouter()

	childRoute.Post("/user/", create)
	childRoute.Get("/user/:id", read)
	childRoute.Put("/user/:id", update)
	childRoute.Delete("/user/:id", delete)

	m.Action(childRoute.Handle)
}
