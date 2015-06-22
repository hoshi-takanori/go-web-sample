package main

import (
	"net/http"

	"github.com/yosssi/ace"
)

var store SessionStore

func main() {
	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	store, err = InitStore(config)
	if err != nil {
		panic(err)
	}

	if config.PublicDir != "" {
		StaticDir(config.PublicDir, http.DefaultServeMux)
	}

	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)

	mux := http.NewServeMux()
	http.HandleFunc("/", AuthHandler(mux))

	if config.PrivateDir != "" {
		StaticDir(config.PrivateDir, mux)
	}

	mux.HandleFunc("/password", PasswordHandler)

	mux.HandleFunc("/", HelloHandler)

	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		panic(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	data := NewData(config.Title)
	username, err := store.GetSession(r)
	if err == nil {
		data["Username"] = username
	}
	ExecTemplate(w, "hello", data)
}

func NewData(title string) map[string]interface{} {
	return map[string]interface{}{
		"Title": title,
	}
}

func ExecTemplate(w http.ResponseWriter, name string, data interface{}) {
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
