package main

import (
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"strconv"
	"encoding/json"
)

func getUserHandler(r render.Render, dbMap *gorp.DbMap, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	logIf(err)

	user, err := dbMap.Get(SecureUser{}, id)
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

	redisConnection.Do("SET", userSession.Id, responseData)

	r.JSON(http.StatusOK, userSession)
}

func logoutHandler(userSession UserSession, r render.Render, redisPool *redis.Pool) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	redisConnection.Do("DEL", userSession.Id)
	r.Error(http.StatusOK)
}
