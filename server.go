package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/piotrkowalczuk/netwars-backend/database"
	"log"
	"net/http"
)

func main() {
	config := ReadConfiguration("config.xml")

	m := martini.New()

	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(render.Renderer())

	redisPool := database.InitializeRedis(config.Redis)
	postgrePool := database.InitPostgre(config.Postgre)

	repositoryManager := NewRepositoryManager(postgrePool, redisPool)

	m.Map(repositoryManager)
	m.Map(redisPool)

	InitRoute(m)

	log.Println("listening on " + config.Server.Host + ":" + config.Server.Port)
	log.Fatalln(http.ListenAndServe(config.Server.Host+":"+config.Server.Port, m))
}

func InitRoute(m *martini.Martini) {
	router := martini.NewRouter()

	CreateForumRoute(router)
	CreateUserRoute(router)
	CreateSearchRoute(router)
	CreateStreamRoute(router)
	//CreateReplayRoute(router)

	m.Action(router.Handle)
}
