package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/piotrkowalczuk/netwars-backend/database"
	"github.com/piotrkowalczuk/netwars-backend/middleware"
	"io"
	"log"
	"testing"
)

var (
	config      *Config
	m           *martini.ClassicMartini
	postgrePool *sql.DB
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

func InitializeEnvironment() {
	config = ReadConfiguration("config_test.xml")

	redisPool := database.InitializeRedis(config.Redis)
	postgrePool := database.InitPostgre(config.Postgre)
	repositoryManager := NewRepositoryManager(postgrePool, redisPool)

	m = martini.Classic()
	m.Use(martini.Recovery())
	m.Use(render.Renderer())
	m.Use(middleware.Sentry(config.Sentry))

	m.Map(repositoryManager)
	m.Map(redisPool)
}

func CreateJSONBody(object interface{}) io.Reader {
	body, err := json.Marshal(object)
	if err != nil {
		log.Println("Unable to marshal object")
	}

	return bytes.NewReader(body)
}
