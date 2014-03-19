package forum

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/binding"
	"github.com/piotrkowalczuk/netwars-backend/user"
)

func CreateRoute(router martini.Router) () {
	router.Get("/forum/:id", getForumHandler)
	router.Get("/forums", getForumsHandler)
	router.Get("/topics/:forumId", getTopicsHandler)
	router.Get("/topic/:id", getTopicHandler)
	router.Get("/posts/:topicId", getPostsHandler)
	router.Post(
		"/post",
		binding.Json(Post{}),
		binding.Form(user.APICredentials{}),
		postPostHandler,
	)
	router.Post(
		"/topic",
		binding.Json(CreateTopicRequest{}),
		binding.Form(user.APICredentials{}),
		postTopicHandler,
	)
}
