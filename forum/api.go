package forum

import (
	"github.com/martini-contrib/render"
	"github.com/piotrkowalczuk/netwars-backend/user"
	"github.com/go-martini/martini"
	"github.com/coopernurse/gorp"
	"github.com/modcloth/sqlutil"
	"github.com/kennygrant/sanitize"
	"database/sql"
	"net/http"
	"strconv"
	"net"
	"time"
	"log"
)

func getForumHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) {
	forumId, _ := strconv.Atoi(params["id"])
	var forum Forum

	err := dbMap.SelectOne(&forum, "SELECT * FROM forum WHERE forum_id = $1", forumId)

	if err != nil {
		r.Error(http.StatusNotFound)
	} else {
		r.JSON(http.StatusOK, &forum)
	}
}

func getForumsHandler(r render.Render, dbMap *gorp.DbMap) {
	var forums []Forum

	_, err := dbMap.Select(&forums, "SELECT * FROM forum")

	if err != nil {
		r.Error(http.StatusNotFound)
	} else {
		r.JSON(http.StatusOK, &forums)
	}
}


func getForumTopicsHandler(r render.Render, req *http.Request, dbMap *gorp.DbMap, params martini.Params) {
	forumId, _ := strconv.Atoi(params["id"])
	queryString := req.URL.Query()

	var limit int
	var offset int
	var topics []Topic

	if limitString := queryString.Get("limit") ; limitString == "" {
		limit = 10
	} else {
		limit, _ = strconv.Atoi(limitString)
	}

	if offsetString := queryString.Get("offset") ; offsetString == "" {
		offset = 0
	} else {
		offset, _ = strconv.Atoi(offsetString)
	}

	_, err := dbMap.Select(
		&topics,
		"SELECT * FROM forum_topic WHERE forum_id = :forumId ORDER BY last_post_date DESC LIMIT :limit OFFSET :offset",
		map[string]interface{}{
			"forumId": forumId,
			"limit": limit,
			"offset": offset,
		},
	)

	if err != nil {
		log.Println(err)
		r.Error(http.StatusNotFound)
	} else {
		r.JSON(http.StatusOK, &topics)
	}
}

func getTopicHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) {
	topicId, _ := strconv.Atoi(params["id"])
	var topic Topic

	err := dbMap.SelectOne(&topic, "SELECT * FROM forum_topic WHERE topic_id = $1", topicId)

	if err != nil {
		log.Println(err)
		r.Error(http.StatusNotFound)
	} else {
		r.JSON(http.StatusOK, &topic)
	}
}

func getTopicPostsHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) {
	topicId, _ := strconv.Atoi(params["id"])
	var posts []Post

	_, err := dbMap.Select(&posts, "SELECT * FROM forum_post WHERE topic_id = $1", topicId)

	if err != nil {
		log.Println(err)
		r.Error(http.StatusNotFound)
	} else {
		for i := range posts {
			posts[i].calculateCreationDiff()
		}
		r.JSON(http.StatusOK, &posts)
	}
}

func postPostHandler(post Post, userSession user.UserSession, req *http.Request, r render.Render, dbMap *gorp.DbMap) {
	var topic Topic

	now := time.Now()

	err := dbMap.SelectOne(&topic, "SELECT * FROM forum_topic WHERE topic_id = $1", post.TopicId)

	post.ChangeAt = &now
	post.CreatedAt = &now
	post.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	post.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	post.AuthorIP.Valid = true
	post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	sanitizedContent := sanitize.HTML(*post.Content)
	post.Content = &sanitizedContent

	err = dbMap.Insert(&post)

	if err != nil {
		log.Println(err)
	}

	topic.NbOfPosts.Int64 += 1
	topic.LastPostDate = &now
	topic.LastPostId = sqlutil.NullInt64{sql.NullInt64{post.Id, true}}
	topic.LastPostAuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	topic.LastPostAuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}

	_, err = dbMap.Update(&topic)

	if err != nil {
		log.Println(err)
		r.Error(http.StatusInternalServerError)
	} else {
		post.calculateCreationDiff()
		r.JSON(http.StatusOK, &post)
	}
}

func patchPostHandler(createPostRequest CreatePostRequest, userSession user.UserSession, params martini.Params, r render.Render, dbMap *gorp.DbMap) {
	if createPostRequest.isValid() {
		postId, _ := strconv.Atoi(params["id"])
		var post Post

		now := time.Now()

		err := dbMap.SelectOne(&post, "SELECT * FROM forum_post WHERE post_id = $1", postId)

		post.ChangeAt = &now
		post.ChangerId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
		post.ChangerName = sqlutil.NullString{sql.NullString{userSession.Name, true}}

		sanitizedContent := sanitize.HTML(*createPostRequest.Content)
		post.Content = &sanitizedContent

		_, err = dbMap.Update(&post)

		if err != nil {
			r.Error(http.StatusInternalServerError)
		} else {
			r.JSON(http.StatusOK, &post)
		}
	} else {
		r.Error(http.StatusBadRequest)
	}
}

func postTopicHandler(createTopicRequest CreateTopicRequest, userSession user.UserSession, req *http.Request, r render.Render, dbMap *gorp.DbMap) {
	if createTopicRequest.isValid() {
		now := time.Now()

		createTopicRequest.Topic.ChangeAt = &now
		createTopicRequest.Topic.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
		createTopicRequest.Topic.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
		createTopicRequest.Topic.LastPostAuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
		createTopicRequest.Topic.LastPostAuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
		createTopicRequest.Topic.NbOfPosts = sqlutil.NullInt64{sql.NullInt64{1, true}}
		createTopicRequest.Topic.NbOfViews = sqlutil.NullInt64{sql.NullInt64{0, true}}

		err := dbMap.Insert(&createTopicRequest.Topic)

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
		} else {
			r.JSON(http.StatusOK, &createTopicRequest.Topic)
		}
	} else {
		r.Error(http.StatusBadRequest)
	}
}


