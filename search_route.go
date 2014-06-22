package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func CreateSearchRoute(router martini.Router) {
	router.Post(
		"/search",
		binding.Json(SearchParams{}),
		binding.Form(APICredentials{}),
		AuthenticationMiddleware(false),
		postSearchHandler,
	)
	router.Get(
		"/search/:id",
		binding.Form(APICredentials{}),
		AuthenticationMiddleware(false),
		getSearchResultsHandler,
	)
}
