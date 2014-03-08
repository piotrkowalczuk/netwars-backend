package server

import (
	"github.com/daviddengcn/go-colortext"
	"github.com/piotrkowalczuk/netwars-backend/database"
	"net/http"
	"log"
)

func Run() {
	database.InitPostgre()

	router := CreateRoute()

	http.Handle("/", router)

	server := &http.Server{
		Addr: "127.0.0.1:3000",
	}

	ct.ChangeColor(ct.Yellow, false, ct.None, false)
	log.Println("Listening on port 3000.")

	ct.ChangeColor(ct.Red, false, ct.None, false)
	log.Fatal(server.ListenAndServe())
}
