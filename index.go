package main

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, user User) {
	FileListHandler(w, user, true)
}

func DailyHandler(w http.ResponseWriter, r *http.Request, user User) {
	FileListHandler(w, user, false)
}

func FileListHandler(w http.ResponseWriter, user User, order bool) {
	year := 0
	if user.year == config.FreshYear {
		year = user.year
	}

	users, err := pgStore.ListUsers(year, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := NewData(config.Title)
	data["Username"] = user.name
	data["Order"] = order

	list := ListFiles(users)
	if order {
		data["Sections"] = MakeSections(list, year)
	} else {
		data["Sections"] = MakeDailySections(list)
	}

	ExecTemplate(w, "index", data)
}
