package main

import (
	"database/sql"
)

type TopicRepository struct {
	db *sql.DB
}

func NewTopicRepository(db *sql.DB) (repository *TopicRepository) {
	repository = &TopicRepository{db}

	return
}

func (tr *TopicRepository) Insert(topic *Topic) (int, error) {
	var id int

	query := `
		INSERT INTO forum_topic (forum_id, topic_name, first_poster, first_poster_name, last_poster, last_poster_name, last_post_id, last_post_date, topic_posts, topic_views, topic_closed, topic_pined, topic_visible_from, topic_visible_to, topic_deleted, change_date, change_user_id, change_ip)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING topic_id
	`
	err := tr.db.QueryRow(
		query,
		&topic.ForumId,
		&topic.Name,
		&topic.AuthorId,
		&topic.AuthorName,
		&topic.LastPostAuthorId,
		&topic.LastPostAuthorName,
		&topic.LastPostId,
		&topic.LastPostDate,
		&topic.NbOfPosts,
		&topic.NbOfViews,
		&topic.IsClosed,
		&topic.IsPinned,
		&topic.VisibleFrom,
		&topic.VisibleTo,
		&topic.IsDeleted,
		&topic.ChangeAt,
		&topic.ChangerId,
		&topic.ChangerIP,
	).Scan(&id)

	return id, err
}

func (tr *TopicRepository) Update(topic *Topic) (sql.Result, error) {
	query := `
		UPDATE forum_topic
		SET
			forum_id = $1,
			topic_name = $2,
			first_poster = $3,
			first_poster_name = $4,
			last_poster = $5,
			last_poster_name = $6,
			last_post_id = $7,
			last_post_date = $8,
			topic_posts = $9,
			topic_views = $10,
			topic_closed = $11,
			topic_pined = $12,
			topic_visible_from = $13,
			topic_visible_to = $14,
			topic_deleted = $15,
			change_date = $16,
			change_user_id = $17,
			change_ip = $18
		WHERE topic_id = $19
	`
	result, err := tr.db.Exec(
		query,
		&topic.ForumId,
		&topic.Name,
		&topic.AuthorId,
		&topic.AuthorName,
		&topic.LastPostAuthorId,
		&topic.LastPostAuthorName,
		&topic.LastPostId,
		&topic.LastPostDate,
		&topic.NbOfPosts,
		&topic.NbOfViews,
		&topic.IsClosed,
		&topic.IsPinned,
		&topic.VisibleFrom,
		&topic.VisibleTo,
		&topic.IsDeleted,
		&topic.ChangeAt,
		&topic.ChangerId,
		&topic.ChangerIP,
		&topic.Id,
	)

	return result, err
}

func (tr *TopicRepository) FindOne(id int64) (error, *Topic) {
	topic := new(Topic)

	err := tr.db.QueryRow(
		"SELECT * FROM forum_topic WHERE topic_id = $1",
		id,
	).Scan(
		&topic.ForumId,
		&topic.Id,
		&topic.Name,
		&topic.AuthorId,
		&topic.AuthorName,
		&topic.LastPostAuthorId,
		&topic.LastPostAuthorName,
		&topic.LastPostId,
		&topic.LastPostDate,
		&topic.NbOfPosts,
		&topic.NbOfViews,
		&topic.IsClosed,
		&topic.IsPinned,
		&topic.VisibleFrom,
		&topic.VisibleTo,
		&topic.IsDeleted,
		&topic.ChangeAt,
		&topic.ChangerId,
		&topic.ChangerIP,
	)

	return err, topic
}

func (tr *TopicRepository) Find(forumId int64, limit int64, offset int64) ([]*Topic, error) {
	topics := []*Topic{}
	var err error

	rows, err := tr.db.Query(
		"SELECT * FROM forum_topic WHERE forum_id = $1 ORDER BY last_post_date DESC LIMIT $2 OFFSET $3",
		forumId,
		limit,
		offset,
	)
	defer rows.Close()

	for rows.Next() {
		topic := new(Topic)
		err = rows.Scan(
			&topic.ForumId,
			&topic.Id,
			&topic.Name,
			&topic.AuthorId,
			&topic.AuthorName,
			&topic.LastPostAuthorId,
			&topic.LastPostAuthorName,
			&topic.LastPostId,
			&topic.LastPostDate,
			&topic.NbOfPosts,
			&topic.NbOfViews,
			&topic.IsClosed,
			&topic.IsPinned,
			&topic.VisibleFrom,
			&topic.VisibleTo,
			&topic.IsDeleted,
			&topic.ChangeAt,
			&topic.ChangerId,
			&topic.ChangerIP,
		)

		topics = append(topics, topic)
	}

	return topics, err
}
