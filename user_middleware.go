package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"net/http"
)

func AuthenticationMiddleware(c martini.Context, apiCredentials APICredentials, res http.ResponseWriter, redisPool *redis.Pool) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.getSessionKey()))
	logIf(err)

	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var userSession UserSession
	json.Unmarshal(userSessionBytes, &userSession)

	if userSession.Token != apiCredentials.Token {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	redisConnection.Do("EXPIRE", userSession.getSessionKey(), SESSION_LIFE_TIME)
	c.Map(userSession)
}
