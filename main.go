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

	CookieName   string
	CookieMaxAge int
	CookieSecure bool

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
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     config.CookieName,
		Path:     "/",
		MaxAge:   -1,
		Secure:   config.CookieSecure,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 302)
}

func hello(w http.ResponseWriter, r *http.Request) {
	template(w, "hello", map[string]string{"Title": "Go Web Sample"})
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
