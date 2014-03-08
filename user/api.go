package user

import (
	"net/http"
	"log"
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"
)

func create(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.Write([]byte(r.URL.Path))
}

func read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	repository := NewUserRepository()
	user := repository.FindOne(id)

	if user == nil {
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
