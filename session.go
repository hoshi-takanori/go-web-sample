package main

import (
	"net/http"
)

type SessionStore interface {
	CheckPassword(username, password string) bool
	ChangePassword(username, password string) error

	StartSession(username string) (string, error)
	ClearSession(r *http.Request)
	GetSession(r *http.Request) (string, error)
}
