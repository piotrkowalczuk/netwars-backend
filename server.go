package main

import (
	"github.com/piotrkowalczuk/netwars-backend/database"
	"github.com/piotrkowalczuk/netwars-backend/user"
	"github.com/piotrkowalczuk/netwars-backend/forum"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"os"
)

func main() {
	config := ReadConfiguration()

	m := martini.New()

	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(render.Renderer())

	redisPool := database.InitializeRedis(config.Redis)

	dbMap := database.InitializeGorp(config.Postgre)
	dbMap.AddTableWithName(user.User{}, "users").SetKeys(true, "user_id")
	dbMap.AddTableWithName(user.SecureUser{}, "users").SetKeys(true, "user_id")
	dbMap.AddTableWithName(forum.Forum{}, "forum").SetKeys(true, "forum_id")
	dbMap.AddTableWithName(forum.Topic{}, "forum_topic").SetKeys(true, "topic_id")
	dbMap.AddTableWithName(forum.Post{}, "forum_post").SetKeys(true, "post_id")

	dbMap.TraceOn("[gorp]", log.New(os.Stdout, "netwars:", log.Lmicroseconds))

	m.Map(dbMap)
	m.Map(redisPool)

	InitRoute(m)

	log.Println("listening on " + config.Server.Host + ":" + config.Server.Port)
	log.Fatalln(http.ListenAndServe(config.Server.Host + ":" + config.Server.Port, m))
}

func InitRoute(m *martini.Martini) () {
	router := martini.NewRouter()

	forum.CreateRoute(router)
	user.CreateRoute(router)

	m.Action(router.Handle)
}

