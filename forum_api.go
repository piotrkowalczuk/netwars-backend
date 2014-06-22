package main

import (
	"database/sql"
	"github.com/go-martini/martini"
	"github.com/kennygrant/sanitize"
	"github.com/martini-contrib/render"
	"github.com/modcloth/sqlutil"
	"net"
	"net/http"
	"github.com/piotrkowalczuk/netwars-backend/service"
	"strconv"
	"time"
)

func getForumHandler(
	r render.Render,
	rm *RepositoryManager,
	params martini.Params,
	sentry *service.Sentry,
) {
	forumId, err := strconv.ParseInt(params["id"], 10, 64)
	sentry.Error(err)

	forum, err := rm.ForumRepository.FindOne(forumId)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &forum)
}

func getForumsHandler(
	r render.Render,
	rm *RepositoryManager,
	sentry *service.Sentry,
) {
	forums, err := rm.ForumRepository.Find()
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &forums)
}

func getForumTopicsHandler(
	rm *RepositoryManager,
	r render.Render,
	req *http.Request,
	userSession UserSession,
	params martini.Params,
	sentry *service.Sentry,
) {
	forumId, err := strconv.ParseInt(params["id"], 10, 64)
	sentry.Error(err)
	queryString := req.URL.Query()

	var topics []*Topic
	var limit int64
	var offset int64

	if limitString := queryString.Get("limit"); limitString == "" {
		limit = int64(10)
	} else {
		limit, err = strconv.ParseInt(limitString, 10, 64)
		sentry.Error(err)
	}

	if offsetString := queryString.Get("offset"); offsetString == "" {
		offset = int64(0)
	} else {
		offset, err = strconv.ParseInt(offsetString, 10, 64)
		sentry.Error(err)
	}

	if userSession.Id == 0 {
		topics, err = rm.TopicRepository.Find(forumId, limit, offset)
		sentry.Error(err)
	} else {
		topics, err = rm.TopicRepository.FindWithUserTopic(forumId, userSession.Id, limit, offset)
		sentry.Error(err)
	}

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &topics)
}

func getTopicHandler(
	r render.Render,
	rm *RepositoryManager,
	params martini.Params,
	userSession UserSession,
	sentry *service.Sentry,
) {
	var topic *Topic
	var err error
	topicId, _ := strconv.ParseInt(params["id"], 10, 64)

	if userSession.Id == 0 {
		topic, err = rm.TopicRepository.FindOne(topicId)
		sentry.Error(err)
	} else {
		topic, err = rm.TopicRepository.FindOneWithUserTopic(topicId, userSession.Id)
		sentry.Error(err)
	}

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &topic)
}

func getTopicPostsHandler(
	r render.Render,
	rm *RepositoryManager,
	params martini.Params,
	userSession UserSession,
	sentry *service.Sentry,
) {
	topicId, _ := strconv.ParseInt(params["id"], 10, 64)

	err, posts := rm.PostRepository.FindByTopicId(topicId)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	if userSession.Id > 0 && len(posts) > 0 {
		userTopic := NewUserTopic(userSession.Id, topicId, len(posts))
		_, err = rm.UserTopicRepository.Upsert(userTopic)
		sentry.Error(err)
	}

	for i := range posts {
		posts[i].calculateCreationDiff()
	}

	r.JSON(http.StatusOK, &posts)
}

func postPostHandler(
	post Post,
	rm *RepositoryManager,
	userSession UserSession,
	req *http.Request,
	r render.Render,
	sentry *service.Sentry,
) {
	now := time.Now()

	topic, err := rm.TopicRepository.FindOne(post.TopicId)
	sentry.Error(err)

	post.ChangeAt = &now
	post.CreatedAt = &now
	post.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	post.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	post.AuthorIP.Valid = true
	post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	sanitizedContent := sanitize.HTML(*post.Content)
	post.Content = &sanitizedContent

	_, err = rm.PostRepository.Insert(&post)
	sentry.Error(err)

	topic.NbOfPosts.Int64 += 1
	topic.LastPostDate = &now
	topic.LastPostId = sqlutil.NullInt64{sql.NullInt64{post.Id, true}}
	topic.LastPostAuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	topic.LastPostAuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}

	_, err = rm.TopicRepository.Update(topic)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	post.calculateCreationDiff()
	r.JSON(http.StatusOK, &post)
}

func patchPostHandler(
	createPostRequest CreatePostRequest,
	userSession UserSession,
	rm *RepositoryManager,
	params martini.Params,
	r render.Render,
	sentry *service.Sentry,
) {
	if !createPostRequest.isValid() {
		r.Error(http.StatusBadRequest)
		return
	}

	postId, err := strconv.ParseInt(params["id"], 10, 64)
	sentry.Error(err)

	now := time.Now()

	err, post := rm.PostRepository.FindOne(postId)
	sentry.Error(err)

	post.ChangeAt = &now
	post.ChangerId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	post.ChangerName = sqlutil.NullString{sql.NullString{userSession.Name, true}}

	sanitizedContent := sanitize.HTML(*createPostRequest.Content)
	post.Content = &sanitizedContent

	_, err = rm.PostRepository.Update(post)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, &post)
}

func postTopicHandler(
	createTopicRequest CreateTopicRequest,
	userSession UserSession,
	rm *RepositoryManager,
	req *http.Request,
	r render.Render,
	sentry *service.Sentry,
) {
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
	createTopicRequest.Topic.LastPostDate = &now
	createTopicRequest.Topic.NbOfPosts = sqlutil.NullInt64{sql.NullInt64{1, true}}
	createTopicRequest.Topic.NbOfViews = sqlutil.NullInt64{sql.NullInt64{0, true}}

	topicId, err := rm.TopicRepository.Insert(&createTopicRequest.Topic)
	sentry.Error(err)

	createTopicRequest.Post.ChangeAt = &now
	createTopicRequest.Post.CreatedAt = &now
	createTopicRequest.Post.TopicId = int64(topicId)
	createTopicRequest.Post.AuthorId = sqlutil.NullInt64{sql.NullInt64{userSession.Id, true}}
	createTopicRequest.Post.AuthorName = sqlutil.NullString{sql.NullString{userSession.Name, true}}
	createTopicRequest.Post.AuthorIP.Valid = true
	createTopicRequest.Post.AuthorIP.String, _, _ = net.SplitHostPort(req.RemoteAddr)

	postId, err := rm.PostRepository.Insert(&createTopicRequest.Post)
	sentry.Error(err)

	createTopicRequest.Topic.Id = int64(topicId)
	createTopicRequest.Topic.LastPostDate = &now
	createTopicRequest.Topic.LastPostId = sqlutil.NullInt64{sql.NullInt64{int64(postId), true}}

	_, err = rm.TopicRepository.Update(&createTopicRequest.Topic)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, &createTopicRequest.Topic)
}
