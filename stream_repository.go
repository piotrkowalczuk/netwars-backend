package main

import (
	"database/sql"
)

type StreamRepository struct {
	db *sql.DB
}

func NewStreamRepository(dbPool *sql.DB) (repository *StreamRepository) {
	repository = &StreamRepository{dbPool}

	return
}

func (usr *StreamRepository) Insert(stream *Stream) (sql.Result, error) {
	query := `
		INSERT INTO user_stream (user_id, streamtype, handle, wol, hots, lol, bw, other)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING user_id
	`
	result, err := usr.db.Exec(
		query,
		&stream.UserId,
		&stream.Type,
		&stream.Identifier,
		&stream.WingsOfLiberty,
		&stream.HeartOfTheSwarm,
		&stream.LeagueOfLegends,
		&stream.BroodWar,
		&stream.Other,
	)

	return result, err
}

func (usr *StreamRepository) Update(stream *Stream) (sql.Result, error) {
	query := `
		UPDATE user_stream
		SET
			streamtype = $1,
			handle = $2,
			wol = $3,
			hots = $4,
			lol = $5,
			bw = $6,
			other = $7,
		WHERE user_id = $8
	`
	result, err := usr.db.Exec(
		query,
		&stream.Type,
		&stream.Identifier,
		&stream.WingsOfLiberty,
		&stream.HeartOfTheSwarm,
		&stream.LeagueOfLegends,
		&stream.BroodWar,
		&stream.Other,
		&stream.UserId,
	)

	return result, err
}

func (usr *StreamRepository) FindOne(id int64) (*Stream, error) {
	stream := new(Stream)

	err := usr.db.QueryRow(
		"SELECT * FROM user_stream WHERE user_id = $1",
		id,
	).Scan(
		&stream.UserId,
		&stream.Type,
		&stream.Identifier,
		&stream.WingsOfLiberty,
		&stream.HeartOfTheSwarm,
		&stream.LeagueOfLegends,
		&stream.BroodWar,
		&stream.Other,
	)

	return stream, err
}

func (usr *StreamRepository) Find(limit int64, offset int64) ([]*Stream, error) {
	streams := []*Stream{}
	var err error

	rows, err := usr.db.Query(
		"SELECT * FROM user_stream ORDER BY last_post_date DESC LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	defer rows.Close()

	for rows.Next() {
		stream := new(Stream)
		err = rows.Scan(
			&stream.UserId,
			&stream.Type,
			&stream.Identifier,
			&stream.WingsOfLiberty,
			&stream.HeartOfTheSwarm,
			&stream.LeagueOfLegends,
			&stream.BroodWar,
			&stream.Other,
		)

		streams = append(streams, stream)
	}

	return streams, err
}

