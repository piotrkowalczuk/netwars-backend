package user

import (
	"github.com/coopernurse/gorp"
	"github.com/codegangsta/martini"
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

func login(credentials Credentials, w http.ResponseWriter, req *http.Request, dbMap *gorp.DbMap) {
	var user User

	/*
	Naive implementation
	 */
	dbMap.SelectOne(
		&user,
		"SELECT * FROM users as u WHERE u.email = $1 AND u.user_pass = $2",
		credentials.Email,
		credentials.Password,
	)

	data, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}
