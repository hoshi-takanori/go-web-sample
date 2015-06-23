package main

import (
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type Section struct {
	Name    string
	Entries []Entry
}

type Entry struct {
	Name string
	Path string
	Date string
	Dcls string
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	data := NewData(config.Title)

	username, err := store.GetSession(r)

	var user *User
	if err == nil {
		user, err = pgStore.GetUser(username)
	}

	var sections []Section
	if err == nil {
		year := 0
		if user.year == config.FreshYear {
			year = user.year
		}
		sections, err = MakeSections(year)
	}

	if err == nil {
		data["Username"] = username
		data["Sections"] = sections
	}

	ExecTemplate(w, "hello", data)
}

func MakeSections(targetYear int) ([]Section, error) {
	users, err := pgStore.ListUsers(targetYear)
	if err != nil {
		return nil, err
	}

	usersMap := map[int][]User{}
	years := []int{}
	for _, user := range users {
		year := user.year
		if targetYear != 0 && year != targetYear {
			year = 0
		}
		list, ok := usersMap[year]
		if ok {
			usersMap[year] = append(list, user)
		} else {
			usersMap[year] = []User{user}
			years = append(years, year)
		}
	}

	sections := []Section{}
	for _, year := range years {
		name := "Staff"
		if year != 0 {
			name = "Fresh"
			if targetYear == 0 {
				name += " " + strconv.Itoa(year)
			}
		}
		sections = append(sections, ListFiles(name, usersMap[year]))
	}

	return sections, nil
}

func ListFiles(name string, users []User) Section {
	entries := []Entry{}
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	for _, user := range users {
		dir := "staff"
		if user.year != 0 {
			dir = strconv.Itoa(user.year)
		}
		file := path.Join(dir, user.name, "diary.html")
		fi, err := os.Stat(path.Join(config.PrivateDir, file))
		date := "-"
		dcls := "date-none"
		if err == nil && fi.Mode().IsRegular() {
			date = fi.ModTime().Format("2006/01/02 15:04:05")
			dcls = DateClass(fi.ModTime(), today)
		} else {
			file = ""
		}
		entries = append(entries, Entry{user.name, file, date, dcls})
	}
	return Section{name, entries}
}

func DateClass(date, today time.Time) string {
	diff := today.Sub(date)
	day := 24 * time.Hour
	if diff <= 0 {
		return "date-today"
	} else if diff <= day {
		return "date-yesterday"
	} else if diff <= 3*day {
		return "date-recent"
	} else if diff <= 7*day {
		return "date-week"
	} else if diff <= 30*day {
		return "date-month"
	} else {
		return "date-old"
	}
}
