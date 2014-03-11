package user

import "github.com/codegangsta/martini"

func CreateRoute(m *martini.Martini) () {

	userRoute := martini.NewRouter()

	userRoute.Post("/user/", create)
	userRoute.Get("/user/:id", read)
	userRoute.Put("/user/:id", update)
	userRoute.Delete("/user/:id", delete)

	m.Action(userRoute.Handle)
}
