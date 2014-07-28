package main

import (
	"log"
	"net/http"
	"time"
)

type listener interface {
	ListenAndServe() error
}

var server listener = &http.Server{
	Addr:         ":3000",
	Handler:      NewServer().handler(),
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

func main() {
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
