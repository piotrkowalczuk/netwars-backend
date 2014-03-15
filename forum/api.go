package forum

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/garyburd/redigo/redis"
	"github.com/coopernurse/gorp"
	"log"
	"net/http"
)

func getForumHandler(r render.Render, redisPool *redis.Pool, dbMap *gorp.DbMap) {
	log.Println("getForumHandler")
}

func getForumsHandler(r render.Render, redisPool *redis.Pool, dbMap *gorp.DbMap) {

	var forums []Forum

	_, err := dbMap.Select(&forums, "SELECT * FROM forum")

	if err != nil {
		r.Error(http.StatusNotFound)
	}

	r.JSON(http.StatusOK, &forums)
}
