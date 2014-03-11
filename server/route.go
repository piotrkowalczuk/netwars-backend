package server

import (
	"github.com/piotrkowalczuk/netwars-backend/user"
	"github.com/codegangsta/martini"
)

func CreateRoute(m *martini.Martini) () {
  	user.CreateRoute(m)
}
