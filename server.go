package main

import (
	"github.com/piotrkowalczuk/netwars-backend/database"
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"net/http"
	"log"
)

func main() {
	config := ReadConfiguration()

	m := martini.New()

	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(render.Renderer())

	redisPool := database.InitializeRedis(config.Redis)
	postgrePool := database.InitPostgre(config.Postgre)

	repositoryManager := NewRepositoryManager(postgrePool)

	m.Map(repositoryManager)
	m.Map(redisPool)

	InitRoute(m)

	log.Println("listening on " + config.Server.Host + ":" + config.Server.Port)
	log.Fatalln(http.ListenAndServe(config.Server.Host + ":" + config.Server.Port, m))
}

func InitRoute(m *martini.Martini) () {
	router := martini.NewRouter()

	CreateForumRoute(router)
	CreateUserRoute(router)
	CreateSearchRoute(router)

	m.Action(router.Handle)
}

