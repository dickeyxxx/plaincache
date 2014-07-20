package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	httpServer := &http.Server{
		Addr:         ":3000",
		Handler:      NewServer().handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(httpServer.ListenAndServe())
}
