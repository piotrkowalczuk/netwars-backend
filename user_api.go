package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"net/http"
	"strconv"
)

func getUserHandler(rm *RepositoryManager, r render.Render, params martini.Params) {
	id, err := strconv.ParseInt(params["id"], 10, 64)
	logIf(err)

	user, err := rm.UserRepository.FindOne(id)
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &user)
}

func registerHandler(userRegistration UserRegistration, r render.Render, rm *RepositoryManager) {
	if !userRegistration.isValid() {
		r.Error(http.StatusBadRequest)
		return
	}

	user := userRegistration.createUser()

	_, err := rm.UserRepository.Insert(user)
	logIf(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.Error(http.StatusOK)
}

func loginHandler(credentials LoginCredentials, r render.Render, redisPool *redis.Pool, rm *RepositoryManager) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	/*
	Naive implementation
	 */
	err, user := rm.UserRepository.FindOneByEmailAndPassword(credentials.Email, credentials.Password)
	logIf(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	userSession := NewUserSession(user)

	responseData, _ := json.Marshal(userSession)

	redisConnection.Do("SET", userSession.getSessionKey(), responseData)
	redisConnection.Do("EXPIRE", userSession.getSessionKey(), SESSION_LIFE_TIME)

	r.JSON(http.StatusOK, userSession)
}

func logoutHandler(userSession UserSession, r render.Render, redisPool *redis.Pool) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	redisConnection.Do("DEL", userSession.getSessionKey())
	r.Error(http.StatusOK)
}


func getOnlineUsersHandler(r render.Render, redisPool *redis.Pool) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	usersOnline := []BasicUser{}
	var keys []interface{}

	keysValues, err := redis.Values(redisConnection.Do("KEYS", "user:*"))
	logIf(err)

	for len(keysValues) > 0 {
		var key interface{}

		keysValues, err = redis.Scan(keysValues, &key)
		keys = append(keys, key)
	}

	values, err := redis.Values(redisConnection.Do("MGET", keys...))
	logIf(err)

	for len(values) > 0 {
		var value []byte
		var user BasicUser

		values, err = redis.Scan(values, &value)

		json.Unmarshal(value, &user)
		usersOnline = append(usersOnline, user)
	}

	r.JSON(http.StatusOK, usersOnline)
}
