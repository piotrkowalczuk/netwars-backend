package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/piotrkowalczuk/netwars-backend/service"
	"net/http"
)

func AuthenticationMiddleware(isOptional bool) martini.Handler {
	return func(
		c martini.Context,
		apiCredentials APICredentials,
		res http.ResponseWriter,
		redisPool *redis.Pool,
		sentry *service.Sentry,
	) {
		redisConnection := redisPool.Get()
		defer redisConnection.Close()

		userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.getSessionKey()))
		sentry.Error(err)

		if err != nil {
			if isOptional {
				var userSession UserSession
				c.Map(userSession)
				return
			} else {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
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
}
