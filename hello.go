package main

import (
	"net/http"
	"os"
	"path"
	"time"
)

var freshList = []string{"fresh1", "fresh2", "fresh3"}
var staffList = []string{"staff1", "staff2", "staff3"}

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
	if err == nil {
		data["Username"] = username
		data["Sections"] = []Section{
			ListFiles("Fresh", "fresh", "diary.html", freshList),
			ListFiles("Staff", "staff", "diary.html", staffList),
		}
	}
	ExecTemplate(w, "hello", data)
}

func ListFiles(name, dir, diary string, users []string) Section {
	entries := []Entry{}
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	for _, user := range users {
		file := path.Join(dir, user, diary)
		fi, err := os.Stat(path.Join(config.PrivateDir, file))
		date := "-"
		dcls := "date-none"
		if err == nil && fi.Mode().IsRegular() {
			date = fi.ModTime().Format("2006/01/02 15:04:05")
			dcls = DateClass(fi.ModTime(), today)
		} else {
			file = ""
		}
		entries = append(entries, Entry{user, file, date, dcls})
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
