package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/yosssi/ace"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	Address   string
	StaticDir string

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

	if config.StaticDir != "" {
		staticDir(config.StaticDir, http.DefaultServeMux)
	}

	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/", hello)

	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		panic(err)
	}
}

func staticDir(dirname string, mux *http.ServeMux) {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.Dir(dirname))
	for _, fi := range fis {
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
			setCookie(w, username, keepLogin)
			http.Redirect(w, r, "/", 302)
			return
		}

		data["Error"] = "Unknown username or password."
	}

	template(w, "login", data)
}

func logout(w http.ResponseWriter, r *http.Request) {
	setCookie(w, "", "")
	http.Redirect(w, r, "/", 302)
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

func hello(w http.ResponseWriter, r *http.Request) {
	data := newData(config.Title)
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err == nil {
		data["Username"] = cookie.Value
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
