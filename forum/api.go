package forum

import (
	"github.com/codegangsta/martini-contrib/render"
	"github.com/piotrkowalczuk/netwars-backend/user"
	"github.com/codegangsta/martini"
	"github.com/coopernurse/gorp"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"net"
	"strconv"
	"time"
	"encoding/json"
)

func getForumHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) {
	forumId, _ := strconv.Atoi(params["id"])

	var forum Forum

	err := dbMap.SelectOne(&forum, "SELECT * FROM forum WHERE forum_id = $1", forumId)

	if err != nil {
		r.Error(http.StatusNotFound)
	}

	r.JSON(http.StatusOK, &forum)
}

func getForumsHandler(r render.Render, dbMap *gorp.DbMap) {
	var forums []Forum

	_, err := dbMap.Select(&forums, "SELECT * FROM forum")

	if err != nil {
		r.Error(http.StatusNotFound)
	}

	r.JSON(http.StatusOK, &forums)
}

func getTopicsHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) string {
	forumId, _ := strconv.Atoi(params["forumId"])
	var topics []Topic

	_, err := dbMap.Select(&topics, "SELECT * FROM forum_topic WHERE forum_id = $1", forumId)

	if err != nil {
		r.Error(http.StatusNotFound)
		return ""
	}

	r.JSON(http.StatusOK, &topics)
	return ""
}

func getTopicHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) string {
	topicId, _ := strconv.Atoi(params["id"])
	var topic Topic

	err := dbMap.SelectOne(&topic, "SELECT * FROM forum_topic WHERE topic_id = $1", topicId)

	if err != nil {
		r.Error(http.StatusNotFound)
		return ""
	}

	r.JSON(http.StatusOK, &topic)
	return ""
}

func getPostsHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) string {
	topicId, _ := strconv.Atoi(params["topicId"])
	var posts []Post

	_, err := dbMap.Select(&posts, "SELECT * FROM forum_post WHERE topic_id = $1", topicId)

	if err != nil {
		r.Error(http.StatusNotFound)
		return ""
	}

	r.JSON(http.StatusOK, &posts)
	return ""
}

func postPostHandler(post Post, apiCredentials user.APICredentials, req *http.Request, r render.Render, redisPool *redis.Pool, dbMap *gorp.DbMap) string {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.Id))
	var userSession user.UserSession
	json.Unmarshal(userSessionBytes, &userSession)

	now := time.Now()

	post.ChangeAt = &now
	post.CreatedAt = &now
	post.AuthorId.Int64 = userSession.Id
	post.AuthorName.String = userSession.Name
	post.AuthorId.Valid = true
	post.AuthorName.Valid = true
	post.AuthorIP.Valid = true
	post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	err = dbMap.Insert(&post)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return ""
	}

	r.JSON(http.StatusOK, &post)
	return ""
}

func postTopicHandler(createTopicRequest CreateTopicRequest, apiCredentials user.APICredentials, req *http.Request, r render.Render, redisPool *redis.Pool, dbMap *gorp.DbMap) string {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.Id))
	var userSession user.UserSession
	json.Unmarshal(userSessionBytes, &userSession)

	now := time.Now()

	createTopicRequest.Topic.ChangeAt = &now
	createTopicRequest.Topic.AuthorId.Int64 = userSession.Id
	createTopicRequest.Topic.AuthorName.String = userSession.Name
	createTopicRequest.Topic.AuthorId.Valid = true
	createTopicRequest.Topic.AuthorName.Valid = true

	err = dbMap.Insert(&createTopicRequest.Topic)

	createTopicRequest.Post.ChangeAt = &now
	createTopicRequest.Post.CreatedAt = &now
	createTopicRequest.Post.TopicId = createTopicRequest.Topic.Id
	createTopicRequest.Post.AuthorId.Int64 = userSession.Id
	createTopicRequest.Post.AuthorName.String = userSession.Name
	createTopicRequest.Post.AuthorId.Valid = true
	createTopicRequest.Post.AuthorName.Valid = true
	createTopicRequest.Post.AuthorIP.Valid = true
	createTopicRequest.Post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	err = dbMap.Insert(&createTopicRequest.Post)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return ""
	}

	r.JSON(http.StatusOK, &createTopicRequest.Topic)
	return ""
}


