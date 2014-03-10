package user

import "github.com/codegangsta/martini"

func CreateRoute(m *martini.Martini) () {

	userRoute := martini.NewRouter()

	userRoute.Post("/user/", create)
	userRoute.Get("/user/{id:[0-9]+}", read)
	userRoute.Put("/user/{id:[0-9]+}", update)
	userRoute.Delete("/user/{id:[0-9]+}", delete)

	m.Action(userRoute.Handle)
}
