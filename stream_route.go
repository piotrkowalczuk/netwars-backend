package main

import (
	"github.com/go-martini/martini"
)

func CreateStreamRoute(router martini.Router) {
	router.Get("/stream/:id", getStreamHandler)
	router.Get("/streams", getStreamsHandler)
}
