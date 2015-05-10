package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yosssi/ace"
)

type Config struct {
	Address   string
	StaticDir string

	SessionCookie http.Cookie

	AceOptions ace.Options
}

var config Config

func main() {
	str, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, &config)
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
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == password {
			setCookie(w, username)
			http.Redirect(w, r, "/", 302)
			return
		}
	}

	template(w, "login", map[string]string{"Title": "Go Web Sample"})
}

func logout(w http.ResponseWriter, r *http.Request) {
	setCookie(w, "")
	http.Redirect(w, r, "/", 302)
}

func setCookie(w http.ResponseWriter, value string) {
	cookie := config.SessionCookie
	cookie.Value = value
	if value == "" {
		cookie.MaxAge = -1
	}
	http.SetCookie(w, &cookie)
}

func hello(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Go Web Sample",
	}
	cookie, err := r.Cookie(config.SessionCookie.Name)
	if err == nil {
		data["Username"] = cookie.Value
	}
	template(w, "hello", data)
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
