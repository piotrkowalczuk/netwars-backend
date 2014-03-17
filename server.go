package main

import (
	"github.com/piotrkowalczuk/netwars-backend/database"
	"github.com/piotrkowalczuk/netwars-backend/user"
	"github.com/piotrkowalczuk/netwars-backend/forum"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini"
	"log"
	"os"
)

func main() {
	m := martini.New()

	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(render.Renderer())

	dbMap := database.InitializeGorp()
	dbMap.AddTableWithName(user.User{}, "users").SetKeys(true, "user_id")
	dbMap.AddTableWithName(forum.Forum{}, "forum").SetKeys(true, "forum_id")
	dbMap.AddTableWithName(forum.Topic{}, "forum_topic").SetKeys(true, "topic_id")
	dbMap.AddTableWithName(forum.Post{}, "forum_post").SetKeys(true, "post_id")

	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "netwars:", log.Lmicroseconds))

	m.Map(dbMap)

	InitRoute(m)
	database.InitializeRedis(m)
	m.Run()
}

func InitRoute(m *martini.Martini) () {
	router := martini.NewRouter()

	forum.CreateRoute(router)
	user.CreateRoute(router)

	m.Action(router.Handle)
}

