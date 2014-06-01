package main

import (
	"database/sql"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) (repository *PostRepository) {
	repository = &PostRepository{db}

	return
}

func (pr *PostRepository) Insert(post *Post) (int, error) {
	var id int

	query := `
		INSERT INTO forum_post (topic_id, user_id, user_name, post_date, post_body, mod_counter, mod_user_id, mod_user_name, mod_date, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING post_id
	`
	err := pr.db.QueryRow(
		query,
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
	).Scan(&id)

	return id, err
}

func (pr *PostRepository) Update(post *Post) (sql.Result, error) {
	query := `
		UPDATE forum_post
		SET
			topic_id = $1,
			user_id = $2,
			user_name = $3,
			post_date = $4,
			post_body = $5,
			mod_counter = $6,
			mod_user_id = $7,
			mod_user_name = $8,
			mod_date = $9,
			ip_address = $10
		WHERE post_id = $11
	`
	result, err := pr.db.Exec(
		query,
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
		&post.Id,
	)

	return result, err
}

func (pr *PostRepository) FindOne(id int64) (error, *Post) {
	post := new(Post)

	err := pr.db.QueryRow(
		"SELECT * FROM forum_post WHERE post_id = $1",
		id,
	).Scan(
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

	PanicIf(err)

	return err, post
}

func (pr *PostRepository) FindByTopicId(topicId int64) (error, []*Post) {
	posts := []*Post{}

	rows, err := pr.db.Query("SELECT * FROM forum_post WHERE topic_id = $1", topicId)
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

	return err, posts
}
