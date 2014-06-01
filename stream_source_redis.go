package main

import ()

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type StreamSourceRedis struct {
	redis *redis.Pool
}

func NewStreamSourceRedis(redisPool *redis.Pool) (source *StreamSourceRedis) {
	source = &StreamSourceRedis{redisPool}

	return
}

func (ssr *StreamSourceRedis) Get(userId int64) (*Stream, error) {
	stream := new(Stream)
	streamKey := "stream:" + strconv.FormatInt(userId, 10)
	redisConnection := ssr.redis.Get()

	defer redisConnection.Close()

	streamBytes, err := redis.Bytes(redisConnection.Do("GET", streamKey))

	if err != nil {
		return nil, err
	}

	json.Unmarshal(streamBytes, stream)

	return stream, nil
}

func (ssr *StreamSourceRedis) Set(stream *Stream) {
	streamKey := "stream:" + strconv.FormatInt(stream.UserId, 10)
	redisConnection := ssr.redis.Get()

	defer redisConnection.Close()

	redisStream, _ := json.Marshal(stream)

	redisConnection.Do("SET", streamKey, redisStream)
	redisConnection.Do("EXPIRE", streamKey, STREAM_LIFE_TIME)
}
