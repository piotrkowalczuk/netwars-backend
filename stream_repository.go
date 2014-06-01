package main

import (
	"database/sql"
)

type StreamRepository struct {
	SourcePostgre *StreamSourcePostgre
	SourceRedis   *StreamSourceRedis
	SourceTwitch  *StreamSourceTwitch
}

func NewStreamRepository(
	sourcePostgre *StreamSourcePostgre,
	sourceRedis *StreamSourceRedis,
	sourceTwitch *StreamSourceTwitch,
) (repository *StreamRepository) {
	repository = &StreamRepository{sourcePostgre, sourceRedis, sourceTwitch}

	return
}

func (sr *StreamRepository) Insert(stream *Stream) (sql.Result, error) {
	return sr.SourcePostgre.Insert(stream)
}

func (sr *StreamRepository) Update(stream *Stream) (sql.Result, error) {
	return sr.SourcePostgre.Update(stream)
}

func (sr *StreamRepository) FindOne(userId int64) (*Stream, error) {
	streamRedis, err := sr.SourceRedis.Get(userId)

	if err == nil {
		return streamRedis, nil
	}

	streamPostgre, err := sr.SourcePostgre.FindOne(userId)

	if err != nil {
		return nil, err
	}

	streamTwitch, err := sr.SourceTwitch.FetchStream(streamPostgre.Identifier)

	if err != nil {
		return nil, err
	}

	if &streamTwitch.Stream != nil {
		streamPostgre.IsOnline = true
		streamPostgre.Current.Id = streamTwitch.Stream.Id
		streamPostgre.Current.Viewers = streamTwitch.Stream.Viewers
		streamPostgre.Current.Broadcaster = streamTwitch.Stream.Broadcaster
		streamPostgre.Current.Preview = streamTwitch.Stream.Preview
		streamPostgre.Current.Game = streamTwitch.Stream.Game
		streamPostgre.Current.Name = streamTwitch.Stream.Name
	}

	sr.SourceRedis.Set(streamPostgre)

	return streamPostgre, nil
}

func (sr *StreamRepository) Find(limit int64, offset int64) ([]*Stream, error) {
	streams := []*Stream{}
	streamsPostgre, err := sr.SourcePostgre.Find(limit, offset)

	for _, stream := range streamsPostgre {
		streamRedis, err := sr.SourceRedis.Get(stream.UserId)

		if err == nil {
			streams = append(streams, streamRedis)
			continue
		}

		streamTwitch, err := sr.SourceTwitch.FetchStream(stream.Identifier)

		if err != nil {
			continue
		}

		if streamTwitch.Stream.Id > 1 {
			stream.IsOnline = true
			stream.Current.Id = streamTwitch.Stream.Id
			stream.Current.Viewers = streamTwitch.Stream.Viewers
			stream.Current.Broadcaster = streamTwitch.Stream.Broadcaster
			stream.Current.Preview = streamTwitch.Stream.Preview
			stream.Current.Game = streamTwitch.Stream.Game
			stream.Current.Name = streamTwitch.Stream.Name
			stream.Current.DisplayName = streamTwitch.Stream.Channel.DisplayName
		}

		sr.SourceRedis.Set(stream)
		streams = append(streams, stream)
	}

	return streams, err
}
