package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"strconv"
	"github.com/piotrkowalczuk/netwars-backend/service"
)

func getUserHandler(
	rm *RepositoryManager,
	r render.Render,
	params martini.Params,
	sentry *service.Sentry,
) {
	id, err := strconv.ParseInt(params["id"], 10, 64)
	sentry.Error(err)

	user, err := rm.UserRepository.FindOne(id)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &user)
}

func registerHandler(
	userRegistration UserRegistration,
	errors binding.Errors,
	r render.Render,
	rm *RepositoryManager,
	sentry *service.Sentry,
) {
	if _, exists := rm.UserRepository.FindBy("email", userRegistration.Email); len(exists) > 0 {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"email"},
			Classification: "UniqueError",
			Message:        "Użytkownik z takim adresem email już istnieje.",
		})
	}

	if _, exists := rm.UserRepository.FindBy("user_name", userRegistration.Name); len(exists) > 0 {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"name"},
			Classification: "UniqueError",
			Message:        "Użytkownik o takiej nazwie już istnieje.",
		})
	}

	if len(errors) > 0 {
		r.JSON(http.StatusBadRequest, errors)
		return
	}

	user := userRegistration.createUser()

	_, err := rm.UserRepository.Insert(user)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusInternalServerError)
		return
	}

	r.Error(http.StatusOK)
}

func loginHandler(
	credentials LoginCredentials,
	r render.Render,
	redisPool *redis.Pool,
	rm *RepositoryManager,
	sentry *service.Sentry,
) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	/*
		Naive implementation
	*/
	err, user := rm.UserRepository.FindOneByEmailAndPassword(credentials.Email, credentials.Password)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	userSession := NewUserSession(user)

	responseData, err := json.Marshal(userSession)
	sentry.Error(err)

	redisConnection.Do("SET", userSession.getSessionKey(), responseData)
	redisConnection.Do("EXPIRE", userSession.getSessionKey(), SESSION_LIFE_TIME)

	r.JSON(http.StatusOK, userSession)
}

func logoutHandler(
	userSession UserSession,
	r render.Render,
	redisPool *redis.Pool,
) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	redisConnection.Do("DEL", userSession.getSessionKey())
	r.Error(http.StatusOK)
}

func getOnlineUsersHandler(
	r render.Render,
	redisPool *redis.Pool,
	sentry *service.Sentry,
) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	usersOnline := make(map[string]BasicUser)
	var keys []interface{}

	keysValues, err := redis.Values(redisConnection.Do("KEYS", "user:*"))
	sentry.Error(err)

	for len(keysValues) > 0 {
		var key interface{}

		keysValues, err = redis.Scan(keysValues, &key)
		sentry.Error(err)
		keys = append(keys, key)
	}

	values, err := redis.Values(redisConnection.Do("MGET", keys...))
	sentry.Error(err)

	for len(values) > 0 {
		var value []byte
		var user BasicUser

		values, err = redis.Scan(values, &value)
		sentry.Error(err)

		json.Unmarshal(value, &user)
		usersOnline[strconv.FormatInt(user.Id, 10)] = user
	}

	r.JSON(http.StatusOK, usersOnline)
}

func getUserStreamHandler(
	rm *RepositoryManager,
	r render.Render,
	userSession UserSession,
	params martini.Params,
	sentry *service.Sentry,
) {
	stream, err := rm.StreamRepository.FindOne(userSession.Id)
	sentry.Error(err)

	if err != nil {
		r.Error(http.StatusNotFound)
		return
	}

	r.JSON(http.StatusOK, &stream)
}

func postUserStreamHandler(
	streamRequest StreamRequest,
	errors binding.Errors,
	userSession UserSession,
	rm *RepositoryManager,
	r render.Render,
	sentry *service.Sentry,
) {
	if len(errors) > 0 {
		r.JSON(http.StatusBadRequest, errors)
		return
	}

	var err error

	streamRequest.UserId = userSession.Id
	streamRequest.Type = int64(1)
	if stream, _ := rm.StreamRepository.SourcePostgre.FindOne(userSession.Id); stream.UserId != 0 {
		_, err = rm.StreamRepository.Update(&streamRequest.Stream)
	} else {
		_, err = rm.StreamRepository.Insert(&streamRequest.Stream)
	}

	if err != nil {
		sentry.Error(err)
		r.Error(http.StatusInternalServerError)
		return
	}

	r.JSON(http.StatusOK, &streamRequest)
}
