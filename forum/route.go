package forum

import (
	"github.com/codegangsta/martini"
)

func CreateRoute(router martini.Router) () {
	router.Get("/forum/:id", getForumHandler)
	router.Get("/forums", getForumsHandler)
}
