package user

import (
	"github.com/garyburd/redigo/redis"
	"net/http"
	"encoding/json"
	"log"
)

func AuthenticationMiddleware(apiCredentials APICredentials, res http.ResponseWriter, redisPool *redis.Pool) string {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	userSessionBytes, err := redis.Bytes(redisConnection.Do("GET", apiCredentials.Id))

	log.Println(apiCredentials)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return "Not Found"
	}

	var userSession UserSession
	json.Unmarshal(userSessionBytes, &userSession)

	if userSession.Token == apiCredentials.Token {
		res.WriteHeader(http.StatusUnauthorized)
		return "Forbindden"
	}

	res.WriteHeader(200)
	return ""
}
