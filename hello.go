package main

import (
	"net/http"
	"os"
	"path"
	"sort"
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

type UserEntry struct {
	user  User
	entry Entry
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	data := NewData(config.Title)

	username, err := store.GetSession(r)
	if err == nil {
		data["Username"] = username

		sections, err := MakeSections(username)
		if err == nil {
			data["Sections"] = sections
		}
	}

	ExecTemplate(w, "hello", data)
}

func MakeSections(username string) ([]Section, error) {
	user, err := pgStore.GetUser(username)
	if err != nil {
		return nil, err
	}

	year := 0
	if user.year == config.FreshYear {
		year = user.year
	}
	users, err := pgStore.ListUsers(year)
	if err != nil {
		return nil, err
	}

	list := ListFiles(users)
	return MakeYearlySections(list, year)
}

func MakeYearlySections(list []UserEntry, targetYear int) ([]Section, error) {
	years := []int{}
	entries := map[int][]Entry{}
	for _, ue := range list {
		year := ue.user.year
		if targetYear != 0 && year != targetYear {
			year = 0
		}

		list, ok := entries[year]
		if ok {
			entries[year] = append(list, ue.entry)
		} else {
			entries[year] = []Entry{ue.entry}
			years = append(years, year)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	sections := []Section{}
	for _, year := range years {
		name := "Staff"
		if year != 0 {
			name = "Fresh"
			if targetYear == 0 {
				name += " " + strconv.Itoa(year)
			}
		}
		sections = append(sections, Section{name, entries[year]})
	}

	return sections, nil
}

func ListFiles(users []User) []UserEntry {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	list := []UserEntry{}
	for _, user := range users {
		dir := "staff"
		if user.staffYear == 0 && user.year != 0 {
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

		list = append(list, UserEntry{user, Entry{user.name, file, date, dcls}})
	}
	return list
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
