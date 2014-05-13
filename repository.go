package main

import (
	"database/sql"
)
type RepositoryManager struct {
	UserRepository *UserRepository
	StreamRepository *StreamRepository
	ForumRepository *ForumRepository
	TopicRepository *TopicRepository
	PostRepository *PostRepository
	SearchRepository *SearchRepository
}

func NewRepositoryManager(postgrePool *sql.DB) (repositoryManager *RepositoryManager) {
	repositoryManager = &RepositoryManager{}

	repositoryManager.UserRepository = NewUserRepository(postgrePool)
	repositoryManager.StreamRepository = NewStreamRepository(postgrePool)
	repositoryManager.ForumRepository = NewForumRepository(postgrePool)
	repositoryManager.TopicRepository = NewTopicRepository(postgrePool)
	repositoryManager.PostRepository = NewPostRepository(postgrePool)
	repositoryManager.SearchRepository = NewSearchRepository(postgrePool)

	return
}
