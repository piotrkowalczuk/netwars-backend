package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"net/http"
	"encoding/json"
)

func AuthenticationMiddleware(c martini.Context, apiCredentials APICredentials, res http.ResponseWriter, redisPool *redis.Pool) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.Id))
	logIf(err)

	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	var userSession UserSession
	json.Unmarshal(userSessionBytes, &userSession)

	if userSession.Token != apiCredentials.Token {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	c.Map(userSession)
}
