package main

import (
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	data := NewData(config.Title)

	username, err := store.GetSession(r)
	if err == nil {
		data["Username"] = username

		sections, err := MakeSections(username, true)
		if err == nil {
			data["Sections"] = sections
		}
	}

	ExecTemplate(w, "hello", data)
}

func MakeSections(username string, yearly bool) ([]Section, error) {
	user, err := pgStore.GetUser(username)
	if err != nil {
		return nil, err
	}

	year := 0
	if user.year == config.FreshYear {
		year = user.year
	}
	users, err := pgStore.ListUsers(year, yearly)
	if err != nil {
		return nil, err
	}

	list := ListFiles(users)
	if yearly {
		return MakeYearlySections(list, year), nil
	} else {
		return MakeDailySections(list), nil
	}
}
