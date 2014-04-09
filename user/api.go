package user

import (
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"strconv"
	"encoding/json"
)

func getUserHandler(w http.ResponseWriter, r *http.Request, dbMap *gorp.DbMap, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	user, err := dbMap.Get(SecureUser{}, id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

func registerHandler(userRegistration UserRegistration, r render.Render, dbMap *gorp.DbMap) {
	if userRegistration.isValid() {
		user := userRegistration.createUser()

		err := dbMap.Insert(user)

		if err != nil {
			r.Error(http.StatusInternalServerError)
		} else {
			r.JSON(http.StatusOK, map[string]interface{}{})
		}
	} else {
		r.Error(http.StatusBadRequest)
	}
}

func loginHandler(credentials LoginCredentials, r render.Render, redisPool *redis.Pool, dbMap *gorp.DbMap) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	var user User

	/*
	Naive implementation
	 */
	err := dbMap.SelectOne(
		&user,
		"SELECT * FROM users as u WHERE u.email = $1 AND u.user_pass = $2",
		credentials.Email,
		credentials.Password,
	)

	if err != nil {
		r.Error(http.StatusNotFound)
	} else {
		userSession := NewUserSession(&user)

		responseData, _ := json.Marshal(userSession)

		redisConnection.Do("SET", userSession.Id, responseData)

		r.JSON(http.StatusOK, userSession)
	}
}

func logoutHandler(userSession UserSession, r render.Render, redisPool *redis.Pool) {
	redisConnection := redisPool.Get()
	defer redisConnection.Close()

	redisConnection.Do("DEL", userSession.Id)
	r.Error(http.StatusOK)
}
