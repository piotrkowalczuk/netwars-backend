package main

import (
	"database/sql"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) (repository *SearchRepository) {
	repository = &SearchRepository{db}

	return
}

func (sr *SearchRepository) Search(text string, searchRange int64, userId int64) (int64, error) {
	var id int

	query := `
		SELECT do_search($1, $2, $3)
	`
	err := sr.db.QueryRow(
		query,
		text,
		searchRange,
		userId,
	).Scan(&id)

	return int64(id), err
}

func (sr *SearchRepository) FetchSearchResultsPosts(id int64, limit int64, offset int64) ([]*Post, error) {
	posts := []*Post{}

	query := `
		SELECT
			fp.post_id,
			fp.topic_id,
			fp.user_id,
			fp.user_name,
			fp.post_date,
			fp.post_body,
			fp.mod_counter,
			fp.mod_user_id,
			fp.mod_user_name,
			fp.mod_date,
			fp.ip_address
		FROM f_result as fr
		LEFT JOIN forum_post as fp ON fp.post_id = fr.fresult_id
		WHERE fr.fresult_type = $1 AND fr.fsearch_id = $2
		ORDER BY fp.post_date DESC
		LIMIT $3
		OFFSET $4
	`
	rows, err := sr.db.Query(query, SEARCH_RANGE_POSTS, id, limit, offset)
	defer rows.Close()

	for rows.Next() {
		post := new(Post)
		err = rows.Scan(
			&post.Id,
			&post.TopicId,
			&post.AuthorId,
			&post.AuthorName,
			&post.CreatedAt,
			&post.Content,
			&post.NbOfChanges,
			&post.ChangerId,
			&post.ChangerName,
			&post.ChangeAt,
			&post.AuthorIP,
		)

		posts = append(posts, post)
	}

	return posts, err
}

func (sr *SearchRepository) FetchSearchResultsTopics(id int64, limit int64, offset int64) ([]*Topic, error) {
	topics := []*Topic{}

	query := `
		SELECT
			ft.forum_id,
			ft.topic_id,
			ft.topic_name,
			ft.first_poster,
			ft.first_poster_name,
			ft.last_poster,
			ft.last_poster_name,
			ft.last_post_id,
			ft.last_post_date,
			ft.topic_posts,
			ft.topic_views,
			ft.topic_closed,
			ft.topic_pined,
			ft.topic_visible_from,
			ft.topic_visible_to,
			ft.topic_deleted,
			ft.change_date,
			ft.change_user_id,
			ft.change_ip
		FROM f_result as fr
		LEFT JOIN forum_topic as ft ON ft.topic_id = fr.fresult_id
		WHERE fr.fresult_type = $1 AND fr.fsearch_id = $2
		ORDER BY ft.last_post_date DESC
		LIMIT $3
		OFFSET $4
	`
	rows, err := sr.db.Query(query, SEARCH_RANGE_TOPICS, id, limit, offset)
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
