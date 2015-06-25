package main

import (
	"net/http"
	"strings"
)

type Router struct {
	entries []RouterEntry
}

type RouterEntry struct {
	method string
	path   string

	handler  func(http.ResponseWriter, *http.Request, User)
	handler2 func(http.ResponseWriter, *http.Request, User, string)
}

func NewRouter() Router {
	return Router{[]RouterEntry{}}
}

func (r *Router) Handle(method, path string,
	handler func(http.ResponseWriter, *http.Request, User)) {
	r.entries = append(r.entries, RouterEntry{method, path, handler, nil})
}

func (r *Router) HandleWithExtra(method, path string,
	handler2 func(http.ResponseWriter, *http.Request, User, string)) {
	r.entries = append(r.entries, RouterEntry{method, path, nil, handler2})
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, err := store.GetSession(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := pgStore.GetUser(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, e := range router.entries {
		if r.Method == e.method {
			if e.handler != nil && r.URL.Path == e.path {
				e.handler(w, r, *user)
				return
			}
			if e.handler2 != nil && strings.HasPrefix(r.URL.Path, e.path) {
				e.handler2(w, r, *user, r.URL.Path[len(e.path):])
				return
			}
		}
	}

	http.Error(w, "Page Not Found", http.StatusNotFound)
}
