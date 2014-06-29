package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func CreateSearchRoute(router martini.Router) {
	router.Post(
		"/search",
		binding.Json(SearchParams{}),
		AuthenticationMiddleware(false),
		postSearchHandler,
	)
	router.Get(
		"/search/:id",
		AuthenticationMiddleware(false),
		getSearchResultsHandler,
	)
}
