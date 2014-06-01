package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

func postSearchHandler(searchParams SearchParams, userSession UserSession, rm *RepositoryManager, r render.Render) {
	var err error

	searchId, err := rm.SearchRepository.Search(searchParams.Text, searchParams.Range, userSession.Id)
	logIf(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, &searchId)
}

func getSearchResultsHandler(rm *RepositoryManager, params martini.Params, r render.Render, req *http.Request) {
	searchId, _ := strconv.ParseInt(params["id"], 10, 64)
	queryString := req.URL.Query()

	var limit int64
	var offset int64

	if limitString := queryString.Get("limit"); limitString == "" {
		limit = int64(10)
	} else {
		limit, _ = strconv.ParseInt(limitString, 10, 64)
	}

	if offsetString := queryString.Get("offset"); offsetString == "" {
		offset = int64(0)
	} else {
		offset, _ = strconv.ParseInt(offsetString, 10, 64)
	}

	posts, err := rm.SearchRepository.FetchSearchResultsPosts(searchId, limit, offset)
	logIf(err)
	topics, err := rm.SearchRepository.FetchSearchResultsTopics(searchId, limit, offset)
	logIf(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, map[string]interface{}{"posts": &posts, "topics": &topics})
}
