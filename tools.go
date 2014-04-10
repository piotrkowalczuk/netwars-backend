package main

import (
	"log"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func logIf(err error) {
	if err != nil {
		log.Println(err)
	}
}
