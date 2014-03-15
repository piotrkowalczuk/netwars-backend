package user

import (
	"github.com/coopernurse/gorp"
	"github.com/codegangsta/martini"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"log"
	"strconv"
	"encoding/json"
)

func create(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte(r.URL.Path))
}

func read(w http.ResponseWriter, r *http.Request, dbMap *gorp.DbMap, params martini.Params) {
	id, err := strconv.Atoi(params["id"])

	user, err := dbMap.Get(User{}, id)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

func update(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte(r.URL.Path))
}

func delete(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte(r.URL.Path))
}

func login(credentials Credentials, w http.ResponseWriter, redisPool *redis.Pool, dbMap *gorp.DbMap, req *http.Request) {
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
		w.WriteHeader(http.StatusNotFound)
	}

	userSession := NewUserSession(&user)

	responseData, _ := json.Marshal(userSession)

	userId, _ := userSession.Id.Value()
	redisConnection.Send("SET", userId, responseData)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(responseData)
}
