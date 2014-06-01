package database

import (
	"github.com/garyburd/redigo/redis"
)

func InitializeRedis(config RedisConfig) (redisPool *redis.Pool) {
	redisPool = &redis.Pool{
		MaxIdle:   3,
		MaxActive: 10, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Address)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	return
}
