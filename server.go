package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Server struct {
	cache Cacher
}

func NewServer() *Server {
	return &Server{NewCache()}
}

func (s *Server) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			s.get(w, r)
		case "POST":
			s.post(w, r)
		case "DELETE":
			s.deleteCache(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	value, ok := s.cache.Read(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := w.Write([]byte(value))
	if err != nil {
		fmt.Printf("error writing response on GET from %s: %s\n", r.RemoteAddr, err)
	}
}

func (s *Server) post(w http.ResponseWriter, r *http.Request) {
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error reading the POST body from %s: %s\n", r.RemoteAddr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.cache.Write(r.URL.Path, string(value))
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write(value)
	if err != nil {
		fmt.Printf("error on writing POST response: %s", err)
	}
}

func (s *Server) deleteCache(w http.ResponseWriter, r *http.Request) {
	s.cache.Delete(r.URL.Path)
}
