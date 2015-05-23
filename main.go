package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/yosssi/ace"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Address    string
	PublicDir  string
	PrivateDir string

	Title string

	DatabaseDriver string
	DatabaseSource string

	SessionCookie http.Cookie

	AceOptions ace.Options
}

var config Config
var db *sql.DB

func main() {
	str, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, &config)
	if err != nil {
		panic(err)
	}

	db, err = sql.Open(config.DatabaseDriver, config.DatabaseSource)
	if err != nil {
		panic(err)
	}

	if config.PublicDir != "" {
		staticDir(config.PublicDir, http.DefaultServeMux)
	}

	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	mux := http.NewServeMux()
	http.HandleFunc("/", authHandler(mux))

	if config.PrivateDir != "" {
		staticDir(config.PrivateDir, mux)
	}

	mux.HandleFunc("/password", password)

	mux.HandleFunc("/", hello)

	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		panic(err)
	}
}

func staticDir(dirname string, mux *http.ServeMux) {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		println(err.Error())
		return
	}

	fileServer := http.FileServer(http.Dir(dirname))
	for _, fi := range fis {
		if fi.Mode()&os.ModeSymlink != 0 {
			fi2, err := os.Stat(path.Join(dirname, fi.Name()))
			if err == nil {
				fi = fi2
			}
		}
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			mux.Handle("/"+fi.Name()+"/", fileServer)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	data := newData(config.Title)

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		keepLogin := r.FormValue("keep-login")
		if checkPassword(username, password) {
			sessionId, err := startSession(username)
			if err == nil {
				setCookie(w, sessionId, keepLogin)
				http.Redirect(w, r, "/", 302)
			}
		}

		data["Error"] = "Login failed."
	}

	template(w, "login", data)
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(r)
	setCookie(w, "", "")
	http.Redirect(w, r, "/", 302)
}

func password(w http.ResponseWriter, r *http.Request) {
	name := "password"
	data := newData("Change Password")

	if r.Method == "POST" {
		username, err := getSession(r)
		current := r.FormValue("current")
		if err == nil && checkPassword(username, current) {
			new1 := r.FormValue("new1")
			new2 := r.FormValue("new2")
			if new1 == new2 && len(new1) >= 6 && new1 != username {
				err := changePassword(username, new1)
				if err == nil {
					name = "login"
					data["Good"] = "Your password has been changed successfully. Please login again."
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

	template(w, name, data)
}

func checkPassword(username, password string) bool {
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

func changePassword(username, password string) error {
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

func setCookie(w http.ResponseWriter, value string, keepLogin string) {
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

func startSession(username string) (string, error) {
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

func clearSession(r *http.Request) {
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err != nil {
		return
	}
	db.Exec("delete from session where id = $1", cookie.Value)
}

func getSession(r *http.Request) (string, error) {
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

func authHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := getSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
		} else {
			handler.ServeHTTP(w, r)
		}
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	data := newData(config.Title)
	username, err := getSession(r)
	if err == nil {
		data["Username"] = username
	}
	template(w, "hello", data)
}

func newData(title string) map[string]interface{} {
	return map[string]interface{}{
		"Title": title,
	}
}

func template(w http.ResponseWriter, name string, data interface{}) {
	tpl, err := ace.Load("base", name, &config.AceOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
