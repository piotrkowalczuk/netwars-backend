package forum

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/piotrkowalczuk/netwars-backend/user"
)

func CreateRoute(router martini.Router) () {
	router.Get("/forum/:id", getForumHandler)
	router.Get("/forum/:id/topics", getForumTopicsHandler)
	router.Get("/forums", getForumsHandler)
	router.Get("/topic/:id", getTopicHandler)
	router.Get("/topic/:id/posts", getTopicPostsHandler)
	router.Post(
		"/post",
		binding.Json(Post{}),
		binding.Form(user.APICredentials{}),
		user.AuthenticationMiddleware,
		postPostHandler,
	)
	router.Post(
		"/topic",
		binding.Json(CreateTopicRequest{}),
		binding.Form(user.APICredentials{}),
		user.AuthenticationMiddleware,
		postTopicHandler,
	)
	router.Patch(
		"/post/:id",
		binding.Json(CreatePostRequest{}),
		binding.Form(user.APICredentials{}),
		user.AuthenticationMiddleware,
		patchPostHandler,
	)
}
