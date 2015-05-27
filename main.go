package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/yosssi/ace"
)

func main() {
	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	err = InitStore(config)
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
		if CheckPassword(username, password) {
			sessionId, err := StartSession(username)
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
	ClearSession(r)
	setCookie(w, "", "")
	http.Redirect(w, r, "/", 302)
}

func password(w http.ResponseWriter, r *http.Request) {
	name := "password"
	data := newData("Change Password")

	if r.Method == "POST" {
		username, err := GetSession(r)
		current := r.FormValue("current")
		if err == nil && CheckPassword(username, current) {
			new1 := r.FormValue("new1")
			new2 := r.FormValue("new2")
			if new1 == new2 && len(new1) >= 6 && new1 != username {
				err := ChangePassword(username, new1)
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

func authHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := GetSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", 302)
		} else {
			handler.ServeHTTP(w, r)
		}
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	data := newData(config.Title)
	username, err := GetSession(r)
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
