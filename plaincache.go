package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":3000")
	if err != nil {
		fmt.Printf("error resolving address: %s\n\n", err)
		os.Exit(1)
	}

	server := NewServer()
	httpServer := &http.Server{
		Addr:         addr.String(),
		Handler:      server.handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("fatal error while serving: %s", err)
		os.Exit(1)
	}
}
