package main

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
)

type RepositoryManager struct {
	UserRepository   *UserRepository
	StreamRepository *StreamRepository
	ForumRepository  *ForumRepository
	TopicRepository  *TopicRepository
	PostRepository   *PostRepository
	SearchRepository *SearchRepository
}

func NewRepositoryManager(postgrePool *sql.DB, redisPool *redis.Pool) (rm *RepositoryManager) {
	rm = &RepositoryManager{}

	rm.UserRepository = NewUserRepository(postgrePool)
	rm.ForumRepository = NewForumRepository(postgrePool)
	rm.TopicRepository = NewTopicRepository(postgrePool)
	rm.PostRepository = NewPostRepository(postgrePool)
	rm.SearchRepository = NewSearchRepository(postgrePool)

	streamSourceRedis := NewStreamSourceRedis(redisPool)
	streamSourcePostgre := NewStreamSourcePostgre(postgrePool)
	streamSourceTwitch := NewStreamSourceTwitch()
	rm.StreamRepository = NewStreamRepository(
		streamSourcePostgre,
		streamSourceRedis,
		streamSourceTwitch,
	)
	return
}
