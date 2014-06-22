package main

import (
	"database/sql"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) (repository *ForumRepository) {
	repository = &ForumRepository{db}

	return
}

func (fr *ForumRepository) FindOne(id int64) (*Forum, error) {
	forum := new(Forum)

	err := fr.db.QueryRow(
		"SELECT * FROM forum WHERE forum_id = $1",
		id,
	).Scan(
		&forum.Id,
		&forum.Name,
		&forum.Description,
		&forum.Order,
		&forum.Type,
		&forum.Topics,
		&forum.ShowTopics,
	)

	return forum, err
}

func (fr *ForumRepository) Find() ([]*Forum, error) {
	var forums []*Forum
	var err error

	rows, err := fr.db.Query("SELECT * FROM forum")
	defer rows.Close()

	for rows.Next() {
		forum := new(Forum)
		err = rows.Scan(
			&forum.Id,
			&forum.Name,
			&forum.Description,
			&forum.Order,
			&forum.Type,
			&forum.Topics,
			&forum.ShowTopics,
		)

		forums = append(forums, forum)
	}

	return forums, err
}
