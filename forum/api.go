package forum

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	"log"
	"net/http"
	"strconv"
)

func getForumHandler(r render.Render, dbMap *gorp.DbMap) {
	log.Println("getForumHandler")
}

func getForumsHandler(r render.Render, dbMap *gorp.DbMap) {

	var forums []Forum

	_, err := dbMap.Select(&forums, "SELECT * FROM forum")

	if err != nil {
		r.Error(http.StatusNotFound)
	}

	r.JSON(http.StatusOK, &forums)
}

func getTopicsHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) {
	forumId, _ := strconv.Atoi(params["forumId"])
	var topics []Topic

	_, err := dbMap.Select(&topics, "SELECT * FROM forum_topic WHERE forum_id = $1", forumId)

	if err != nil {
		r.Error(http.StatusNotFound)
	}

	r.JSON(http.StatusOK, &topics)
}
