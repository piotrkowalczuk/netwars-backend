package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func CreateForumRoute(router martini.Router) () {
	router.Get("/forum/:id", getForumHandler)
	router.Get("/forum/:id/topics", getForumTopicsHandler)
	router.Get("/forums", getForumsHandler)
	router.Get("/topic/:id", getTopicHandler)
	router.Get("/topic/:id/posts", getTopicPostsHandler)
	router.Post(
		"/post",
		binding.Json(Post{}),
		binding.Form(APICredentials{}),
		AuthenticationMiddleware,
		postPostHandler,
	)
	router.Post(
		"/topic",
		binding.Json(CreateTopicRequest{}),
		binding.Form(APICredentials{}),
		AuthenticationMiddleware,
		postTopicHandler,
	)
	router.Patch(
		"/post/:id",
		binding.Json(CreatePostRequest{}),
		binding.Form(APICredentials{}),
		AuthenticationMiddleware,
		patchPostHandler,
	)
}
