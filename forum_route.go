package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func CreateForumRoute(router martini.Router) {
	router.Get("/forum/:id", getForumHandler)
	router.Get(
		"/forum/:id/topics",
		binding.Form(APICredentials{}),
		AuthenticationMiddleware(true),
		getForumTopicsHandler,
	)
	router.Get("/forums", getForumsHandler)
	router.Get(
		"/topic/:id",
		binding.Form(APICredentials{}),
		AuthenticationMiddleware(true),
		getTopicHandler,
	)
	router.Get(
		"/topic/:id/posts",
		binding.Form(APICredentials{}),
		AuthenticationMiddleware(true),
		getTopicPostsHandler,
	)
	router.Post(
		"/post",
		binding.Json(Post{}),
		binding.Form(APICredentials{}),
		AuthenticationMiddleware(false),
		postPostHandler,
	)
	router.Post(
		"/topic",
		binding.Json(CreateTopicRequest{}),
		binding.Form(APICredentials{}),
		AuthenticationMiddleware(false),
		postTopicHandler,
	)
	router.Patch(
		"/post/:id",
		binding.Json(CreatePostRequest{}),
		binding.Form(APICredentials{}),
		AuthenticationMiddleware(false),
		patchPostHandler,
	)
}
