package main

import (
	"net/http"
	"time"
)

type SessionStore interface {
	CheckPassword(username, password string) bool
	ChangePassword(username, password string) error

	StartSession(username string) (string, error)
	ClearSession(r *http.Request)
	GetSession(r *http.Request) (string, error)
}

func AuthHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := store.GetSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
		} else {
			handler.ServeHTTP(w, r)
		}
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := NewData(config.Title)

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		keepLogin := r.FormValue("keep-login")
		if store.CheckPassword(username, password) {
			sessionId, err := store.StartSession(username)
			if err == nil {
				SetCookie(w, sessionId, keepLogin)
				http.Redirect(w, r, "/", 302)
			}
		}

		data["Error"] = "Login failed."
	}

	ExecTemplate(w, "login", data)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	store.ClearSession(r)
	SetCookie(w, "", "")
	http.Redirect(w, r, "/", 302)
}

func PasswordHandler(w http.ResponseWriter, r *http.Request) {
	name := "password"
	data := NewData("Change Password")

	if r.Method == "POST" {
		username, err := store.GetSession(r)
		current := r.FormValue("current")
		if err == nil && store.CheckPassword(username, current) {
			new1 := r.FormValue("new1")
			new2 := r.FormValue("new2")
			if new1 == new2 && len(new1) >= 6 && new1 != username {
				err := store.ChangePassword(username, new1)
				if err == nil {
					name = "login"
					data["Good"] = "Your password has been changed successfully. " +
						"Please login again."
				} else {
					data["Error"] = "Failed to change your password."
				}
			} else {
				data["Error"] = "Bad new password."
			}
		} else {
			data["Error"] = "Bad current password."
		}
	}

	ExecTemplate(w, name, data)
}

func SetCookie(w http.ResponseWriter, value string, keepLogin string) {
	cookie := config.SessionCookie
	cookie.Value = value
	if value == "" {
		cookie.MaxAge = -1
	} else if keepLogin == "" {
		cookie.MaxAge = 0
	} else if cookie.MaxAge > 0 {
		d := time.Duration(cookie.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	}
	http.SetCookie(w, &cookie)
}
