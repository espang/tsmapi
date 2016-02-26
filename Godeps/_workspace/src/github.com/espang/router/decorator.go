package router

import "net/http"

type Server interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ServerFunc func(http.ResponseWriter, *http.Request)

func (s ServerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s(w, r)
}

type Decorator func(Server) Server
