package forum

import (
	"github.com/martini-contrib/render"
	"github.com/piotrkowalczuk/netwars-backend/user"
	"github.com/go-martini/martini"
	"github.com/coopernurse/gorp"
	"github.com/garyburd/redigo/redis"
	"github.com/modcloth/sqlutil"
	"database/sql"
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

	_, err := dbMap.Select(&topics, "SELECT * FROM forum_topic WHERE forum_id = $1 ORDER BY last_post_date DESC", forumId)

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

	var userSession user.UserSession
	var topic Topic

	now := time.Now()

	userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.Id))
	json.Unmarshal(userSessionBytes, &userSession)

	err = dbMap.SelectOne(&topic, "SELECT * FROM forum_topic WHERE topic_id = $1", post.TopicId)

	post.ChangeAt = &now
	post.CreatedAt = &now
	post.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	post.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	post.AuthorIP.Valid = true
	post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	err = dbMap.Insert(&post)

	topic.NbOfPosts.Int64 += 1
	topic.LastPostDate = &now
	topic.LastPostId = sqlutil.NullInt64{sql.NullInt64{post.Id, true}}
	topic.LastPostAuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	topic.LastPostAuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}

	_, err = dbMap.Update(&topic)

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
	createTopicRequest.Topic.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	createTopicRequest.Topic.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	createTopicRequest.Topic.LastPostAuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	createTopicRequest.Topic.LastPostAuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	createTopicRequest.Topic.NbOfPosts = sqlutil.NullInt64{sql.NullInt64{1, true}}
	createTopicRequest.Topic.NbOfViews = sqlutil.NullInt64{sql.NullInt64{0, true}}

	err = dbMap.Insert(&createTopicRequest.Topic)

	createTopicRequest.Post.ChangeAt = &now
	createTopicRequest.Post.CreatedAt = &now
	createTopicRequest.Post.TopicId = createTopicRequest.Topic.Id
	createTopicRequest.Post.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	createTopicRequest.Post.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	createTopicRequest.Post.AuthorIP.Valid = true
	createTopicRequest.Post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	err = dbMap.Insert(&createTopicRequest.Post)

	createTopicRequest.Topic.LastPostDate = &now
	createTopicRequest.Topic.LastPostId.Int64 = createTopicRequest.Post.Id

	_, err = dbMap.Update(&createTopicRequest.Topic)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return ""
	}

	r.JSON(http.StatusOK, &createTopicRequest.Topic)
	return ""
}


