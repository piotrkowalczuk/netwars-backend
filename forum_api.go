package main

import (
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	"github.com/modcloth/sqlutil"
	"github.com/kennygrant/sanitize"
	"database/sql"
	"net/http"
	"strconv"
	"net"
	"time"
	"log"
)

func getForumHandler(r render.Render, rm *RepositoryManager, params martini.Params) {
	forumId, _ := strconv.ParseInt(params["id"], 10, 64)

	err, forum := rm.ForumRepository.FindOne(forumId)
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &forum)
}

func getForumsHandler(r render.Render, rm *RepositoryManager) {
	err, forums := rm.ForumRepository.Find()
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &forums)
}


func getForumTopicsHandler(rm *RepositoryManager, r render.Render, req *http.Request, params martini.Params) {
	forumId, _ := strconv.ParseInt(params["id"], 10, 64)
	queryString := req.URL.Query()

	var limit int64
	var offset int64

	if limitString := queryString.Get("limit") ; limitString == "" {
		limit = int64(10)
	} else {
		limit, _ = strconv.ParseInt(limitString, 10, 64)
	}

	if offsetString := queryString.Get("offset") ; offsetString == "" {
		offset = int64(0)
	} else {
		offset, _ = strconv.ParseInt(offsetString, 10, 64)
	}

	err, topics := rm.TopicRepository.Find(forumId, limit, offset)
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &topics)
}

func getTopicHandler(r render.Render, rm *RepositoryManager, params martini.Params) {
	topicId, _ := strconv.ParseInt(params["id"], 10, 64)

	err, topic := rm.TopicRepository.FindOne(topicId)
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &topic)
}

func getTopicPostsHandler(r render.Render, rm *RepositoryManager, params martini.Params) {
	topicId, _ := strconv.ParseInt(params["id"], 10, 64)

	err, posts := rm.PostRepository.FindByTopicId(topicId)
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	for i := range posts {
		posts[i].calculateCreationDiff()
	}

	r.JSON(http.StatusOK, &posts)
}

func postPostHandler(post Post, rm *RepositoryManager, userSession UserSession, req *http.Request, r render.Render) {
	now := time.Now()

	err, topic := rm.TopicRepository.FindOne(post.TopicId)
	logIf(err)

	post.ChangeAt = &now
	post.CreatedAt = &now
	post.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	post.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	post.AuthorIP.Valid = true
	post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	sanitizedContent := sanitize.HTML(*post.Content)
	post.Content = &sanitizedContent

	_, err = rm.PostRepository.Insert(&post)
	logIf(err)

	topic.NbOfPosts.Int64 += 1
	topic.LastPostDate = &now
	topic.LastPostId = sqlutil.NullInt64{sql.NullInt64{post.Id, true}}
	topic.LastPostAuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	topic.LastPostAuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}

	_, err = rm.TopicRepository.Update(topic)
	logIf(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	post.calculateCreationDiff()
	r.JSON(http.StatusOK, &post)
}

func patchPostHandler(createPostRequest CreatePostRequest, userSession UserSession, rm *RepositoryManager, params martini.Params, r render.Render) {
	if !createPostRequest.isValid() {
		r.Error(http.StatusBadRequest)
		return
	}

	postId, _ := strconv.ParseInt(params["id"], 10, 64)

	now := time.Now()

	err, post := rm.PostRepository.FindOne(postId)
	logIf(err)

	post.ChangeAt = &now
	post.ChangerId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	post.ChangerName = sqlutil.NullString{sql.NullString{userSession.Name, true}}

	sanitizedContent := sanitize.HTML(*createPostRequest.Content)
	post.Content = &sanitizedContent

	_, err = rm.PostRepository.Update(post)
	logIf(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, &post)
}

func postTopicHandler(createTopicRequest CreateTopicRequest, userSession UserSession, rm *RepositoryManager, req *http.Request, r render.Render) {
	if !createTopicRequest.isValid() {
		r.Error(http.StatusBadRequest)
		return
	}

	var err error
	now := time.Now()

	createTopicRequest.Topic.ChangeAt = &now
	createTopicRequest.Topic.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	createTopicRequest.Topic.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	createTopicRequest.Topic.LastPostAuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	createTopicRequest.Topic.LastPostAuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	createTopicRequest.Topic.NbOfPosts = sqlutil.NullInt64{sql.NullInt64{1, true}}
	createTopicRequest.Topic.NbOfViews = sqlutil.NullInt64{sql.NullInt64{0, true}}

	topicId, err := rm.TopicRepository.Insert(&createTopicRequest.Topic)

	createTopicRequest.Post.ChangeAt = &now
	createTopicRequest.Post.CreatedAt = &now
	createTopicRequest.Post.TopicId = int64(topicId)
	createTopicRequest.Post.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	createTopicRequest.Post.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	createTopicRequest.Post.AuthorIP.Valid = true
	createTopicRequest.Post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	postId, err := rm.PostRepository.Insert(&createTopicRequest.Post)
	logIf(err)

	createTopicRequest.Topic.Id = int64(topicId)
	createTopicRequest.Topic.LastPostDate = &now
	createTopicRequest.Topic.LastPostId = sqlutil.NullInt64{sql.NullInt64{int64(postId), true}}

	result, err := rm.TopicRepository.Update(&createTopicRequest.Topic)
	log.Println(result)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, &createTopicRequest.Topic)
}


