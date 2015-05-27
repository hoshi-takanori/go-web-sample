// +build postgres

package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type PostgresStore struct {
	db *sql.DB
}

func InitStore(config Config) (SessionStore, error) {
	db, err := sql.Open(config.DatabaseDriver, config.DatabaseSource)
	return PostgresStore{db}, err
}

func (s PostgresStore) CheckPassword(username, password string) bool {
	if username == "" || password == "" {
		return false
	}

	var hash string
	err := s.db.QueryRow("select password from users where name = $1", username).Scan(&hash)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s PostgresStore) ChangePassword(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("update users set password = $1 where name = $2", string(hash), username)
	if err != nil {
		return err
	}

	s.db.Exec("delete from session where user_name = $1", username)

	return nil
}

func (s PostgresStore) StartSession(username string) (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	sid := strings.TrimRight(base64.URLEncoding.EncodeToString(buf), "=")
	_, err = s.db.Exec("insert into session values ($1, $2)", sid, username)
	if err != nil {
		return "", err
	}

	return sid, nil
}

func (s PostgresStore) ClearSession(r *http.Request) {
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err != nil {
		return
	}

	s.db.Exec("delete from session where id = $1", cookie.Value)
}

func (s PostgresStore) GetSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err != nil {
		return "", err
	}

	var username string
	err = s.db.QueryRow("select user_name from session where id = $1", cookie.Value).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}
