package main

import (
	"github.com/piotrkowalczuk/netwars-backend/database"
	"github.com/piotrkowalczuk/netwars-backend/user"
	"github.com/codegangsta/martini"
	"log"
	"os"
)

func main() {
	m := martini.New()

	m.Use(martini.Logger())
	m.Use(martini.Recovery())

	dbMap := database.InitializeGorp()
	dbMap.AddTableWithName(user.User{}, "users").SetKeys(true, "user_id")
	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "netwars:", log.Lmicroseconds))

	m.Map(dbMap)

	InitRoute(m)

	m.Run()
}

func InitRoute(m *martini.Martini) () {
	user.CreateRoute(m)
}

