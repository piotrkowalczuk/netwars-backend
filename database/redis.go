package database

import (
	"github.com/garyburd/redigo/redis"
	"github.com/codegangsta/martini"
)
func InitializeRedis(m *martini.Martini) {
	redisPool := &redis.Pool{
		MaxIdle: 3,
		MaxActive: 10, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	m.Map(redisPool)
}
