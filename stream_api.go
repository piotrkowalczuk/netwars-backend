package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
)

func getStreamHandler(rm *RepositoryManager, r render.Render, params martini.Params) {
	id, err := strconv.ParseInt(params["id"], 10, 64)
	logIf(err)

	stream, err := rm.StreamRepository.FindOne(id)
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &stream)
}

func getStreamsHandler(rm *RepositoryManager, r render.Render) {
	limit := int64(100)
	offset := int64(0)
	streams, err := rm.StreamRepository.Find(limit, offset)
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &streams)
}
