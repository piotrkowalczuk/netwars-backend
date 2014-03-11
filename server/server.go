package server

import (
	"github.com/codegangsta/martini"
)

func Run() {
	m := martini.New()

	m.Use(martini.Logger())
	m.Use(martini.Recovery())

	CreateRoute(m)

	m.Run()
}
