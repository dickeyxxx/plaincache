package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var m sync.RWMutex

var cache map[string]string

func main() {
	addr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	if err != nil {
		fmt.Printf("error resolving address: %s\n\n", err)
		os.Exit(1)
	}

	cache = make(map[string]string)

	server := &http.Server{
		Addr:         addr.String(),
		Handler:      handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("fatal error while serving: %s", err)
		os.Exit(1)
	}
}

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			get(w, r)
		case "POST":
			post(w, r)
		case "DELETE":
			deleteCache(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	m.RLock()
	value, ok := cache[r.URL.Path]
	m.RUnlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := w.Write([]byte(value))
	if err != nil {
		fmt.Printf("error writing response on GET from %s: %s\n", r.RemoteAddr, err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error reading the POST body from %s: %s\n", r.RemoteAddr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m.Lock()
	cache[r.URL.Path] = string(value)
	m.Unlock()
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write(value)
	if err != nil {
		fmt.Printf("error on writing POST response: %s", err)
	}
}

func deleteCache(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	delete(cache, r.URL.Path)
	m.Unlock()
}
