package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func CreateForumRoute(router martini.Router) {
	router.Get("/forum/:id", getForumHandler)
	router.Get(
		"/forum/:id/topics",
		AuthenticationMiddleware(true),
		getForumTopicsHandler,
	)
	router.Get("/forums", getForumsHandler)
	router.Get(
		"/topic/:id",
		AuthenticationMiddleware(true),
		getTopicHandler,
	)
	router.Get(
		"/topic/:id/posts",
		AuthenticationMiddleware(true),
		getTopicPostsHandler,
	)
	router.Post(
		"/post",
		binding.Json(Post{}),
		AuthenticationMiddleware(false),
		postPostHandler,
	)
	router.Post(
		"/topic",
		binding.Json(CreateTopicRequest{}),
		AuthenticationMiddleware(false),
		postTopicHandler,
	)
	router.Patch(
		"/post/:id",
		binding.Json(CreatePostRequest{}),
		AuthenticationMiddleware(false),
		patchPostHandler,
	)
}
