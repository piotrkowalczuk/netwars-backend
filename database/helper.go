package database

import "log"

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
