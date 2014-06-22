package main

import (
	"database/sql"
)

type UserTopicRepository struct {
	db *sql.DB
}

func NewUserTopicRepository(dbPool *sql.DB) (repository *UserTopicRepository) {
	repository = &UserTopicRepository{dbPool}

	return
}

func (utr *UserTopicRepository) Insert(userTopic *UserTopic) (sql.Result, error) {
	query := `
        INSERT INTO user_topic (user_id, topic_id, observed, post_seen)
        VALUES ($1, $2, $3, $4)
    `
	result, err := utr.db.Exec(
		query,
		&userTopic.UserId,
		&userTopic.TopicId,
		&userTopic.Observed,
		&userTopic.PostSeen,
	)

	return result, err
}

func (utr *UserTopicRepository) Update(userTopic *UserTopic) (sql.Result, error) {
	query := `
        UPDATE user_topic
        SET user_id = $1,
            topic_id = $2,
            observed = $3,
            post_seen = $4
        WHERE user_id = $5 AND topic_id = $6
		RETURNING user_id
    `
	result, err := utr.db.Exec(
		query,
		&userTopic.UserId,
		&userTopic.TopicId,
		&userTopic.Observed,
		&userTopic.PostSeen,
		&userTopic.UserId,
		&userTopic.TopicId,
	)

	return result, err
}

func (utr *UserTopicRepository) Upsert(userTopic *UserTopic) (sql.Result, error) {
	updateResult, err := utr.Update(userTopic)

	if nbOfRows, _ := updateResult.RowsAffected(); nbOfRows < 1 {
		insertResult, err := utr.Insert(userTopic)

		return insertResult, err
	}

	return updateResult, err
}

func (utr *UserTopicRepository) FindOne(userId int64, topicId int64) (*UserTopic, error) {
	userTopic := new(UserTopic)

	err := utr.db.QueryRow(
		"SELECT * FROM user_topic as ut WHERE ut.user_id = $1 AND ut.topic_id = $2",
		userId,
		topicId,
	).Scan(
		&userTopic.UserId,
		&userTopic.TopicId,
		&userTopic.Observed,
		&userTopic.PostSeen,
	)

	return userTopic, err
}
