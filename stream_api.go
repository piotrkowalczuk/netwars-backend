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

func postStreamHandler(streamRequest StreamRequest, userSession UserSession, rm *RepositoryManager, req *http.Request, r render.Render) {
	if !streamRequest.isValid() {
		r.Error(http.StatusBadRequest)
		return
	}

	var err error

	streamRequest.UserId = userSession.Id


	_, err = rm.StreamRepository.Insert(&streamRequest.Stream)

	if err != nil {
	r.Error(http.StatusInternalServerError)
	return
	}

	r.JSON(http.StatusOK, &streamRequest)
}
