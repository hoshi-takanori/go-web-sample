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

var db *sql.DB

func InitStore(config Config) error {
	var err error
	db, err = sql.Open(config.DatabaseDriver, config.DatabaseSource)
	return err
}

func CheckPassword(username, password string) bool {
	if username == "" || password == "" {
		return false
	}

	var hash string
	err := db.QueryRow("select password from users where name = $1", username).Scan(&hash)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ChangePassword(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("update users set password = $1 where name = $2", string(hash), username)
	if err != nil {
		return err
	}

	db.Exec("delete from session where user_name = $1", username)

	return nil
}

func StartSession(username string) (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	sid := strings.TrimRight(base64.URLEncoding.EncodeToString(buf), "=")
	_, err = db.Exec("insert into session values ($1, $2)", sid, username)
	if err != nil {
		return "", err
	}
	return sid, nil
}

func ClearSession(r *http.Request) {
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err != nil {
		return
	}
	db.Exec("delete from session where id = $1", cookie.Value)
}

func GetSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err != nil {
		return "", err
	}
	var username string
	err = db.QueryRow("select user_name from session where id = $1", cookie.Value).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
