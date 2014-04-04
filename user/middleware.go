package user

import (
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"net/http"
	"encoding/json"
	"log"
)

func AuthenticationMiddleware(c martini.Context, apiCredentials APICredentials, res http.ResponseWriter, redisPool *redis.Pool) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.Id))

	if err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusNotFound)
	}

	var userSession UserSession
	json.Unmarshal(userSessionBytes, &userSession)

	if userSession.Token != apiCredentials.Token {
		res.WriteHeader(http.StatusUnauthorized)
	}

	c.Map(userSession)
}
