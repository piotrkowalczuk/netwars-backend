package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"github.com/piotrkowalczuk/netwars-backend/service"
)

func postSearchHandler(
	searchParams SearchParams,
	userSession UserSession,
	rm *RepositoryManager,
	r render.Render,
	sentry *service.Sentry,
) {
	var err error

	searchId, err := rm.SearchRepository.Search(searchParams.Text, searchParams.Range, userSession.Id)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, &searchId)
}

func getSearchResultsHandler(
	rm *RepositoryManager,
	params martini.Params,
	r render.Render,
	req *http.Request,
	sentry *service.Sentry,
) {
	var limit int64
	var offset int64
	var err error

	searchId, err := strconv.ParseInt(params["id"], 10, 64)
	sentry.Error(err)
	queryString := req.URL.Query()

	if limitString := queryString.Get("limit"); limitString == "" {
		limit = int64(10)
	} else {
		limit, err = strconv.ParseInt(limitString, 10, 64)
		sentry.Error(err)
	}

	if offsetString := queryString.Get("offset"); offsetString == "" {
		offset = int64(0)
	} else {
		offset, err = strconv.ParseInt(offsetString, 10, 64)
		sentry.Error(err)
	}

	posts, err := rm.SearchRepository.FetchSearchResultsPosts(searchId, limit, offset)
	sentry.Error(err)
	topics, err := rm.SearchRepository.FetchSearchResultsTopics(searchId, limit, offset)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{"posts": &posts, "topics": &topics})
}
