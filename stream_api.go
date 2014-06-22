package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"github.com/piotrkowalczuk/netwars-backend/service"
	"strconv"
)

func getStreamHandler(
	rm *RepositoryManager,
	r render.Render,
	params martini.Params,
	sentry *service.Sentry,
) {
	id, err := strconv.ParseInt(params["id"], 10, 64)
	sentry.Error(err)

	stream, err := rm.StreamRepository.FindOne(id)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &stream)
}

func getStreamsHandler(
	rm *RepositoryManager,
	r render.Render,
	sentry *service.Sentry,
) {
	limit := int64(100)
	offset := int64(0)
	streams, err := rm.StreamRepository.Find(limit, offset)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &streams)
}
