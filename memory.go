// +build memory

package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type MemoryStore struct {
	users   map[string]string
	session map[string]string
}

type MemoryError struct {
	msg string
}

func (e MemoryError) Error() string {
	return e.msg
}

func InitStore(config Config) (SessionStore, error) {
	return MemoryStore{
		users:   map[string]string{},
		session: map[string]string{},
	}, nil
}

func (s MemoryStore) CheckPassword(username, password string) bool {
	if username == "" || password == "" {
		return false
	}

	hash, ok := s.users[username]
	if !ok {
		err := s.ChangePassword(username, password)
		return err == nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s MemoryStore) ChangePassword(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.users[username] = string(hash)

	for sid, u := range s.session {
		if u == username {
			delete(s.session, sid)
		}
	}

	return nil
}

func (s MemoryStore) StartSession(username string) (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	sid := strings.TrimRight(base64.URLEncoding.EncodeToString(buf), "=")
	s.session[sid] = username

	return sid, nil
}

func (s MemoryStore) ClearSession(r *http.Request) {
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err != nil {
		return
	}

	delete(s.session, cookie.Value)
}

func (s MemoryStore) GetSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err != nil {
		return "", err
	}

	username, ok := s.session[cookie.Value]
	if !ok {
		return "", MemoryError{"session not found"}
	}

	return username, nil
}
