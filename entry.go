package main

import (
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
	Year string
	Path string
	Date string
	Dcls string
}

type UserEntry struct {
	user  User
	path  string
	mtime time.Time
}

func NewEntry(ue UserEntry, today time.Time) Entry {
	date := "-"
	dcls := "date-none"
	if !ue.mtime.IsZero() {
		date = ue.mtime.Format("2006/01/02 15:04:05")
		dcls = DateClass(ue.mtime, today)
	}
	return Entry{ue.user.name, "", ue.path, date, dcls}
}

func MakeSections(list []UserEntry, targetYear int) []Section {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, now.Location())

	years := []int{}
	entries := map[int][]Entry{}
	for _, ue := range list {
		year := ue.user.year
		if targetYear != 0 && year != targetYear {
			year = 0
		}

		list, ok := entries[year]
		if !ok {
			list = []Entry{}
			years = append(years, year)
		}
		entries[year] = append(list, NewEntry(ue, today))
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

	return sections
}

type ByDate []UserEntry

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].mtime.After(a[j].mtime) }

func NewDailyEntry(ue UserEntry) Entry {
	year := "staff"
	if ue.user.year != 0 && ue.user.staffYear == 0 {
		year = strconv.Itoa(ue.user.year)
	}
	date := "-"
	if !ue.mtime.IsZero() {
		date = ue.mtime.Format("15:04:05")
	}
	return Entry{ue.user.name, year, ue.path, date, "date-none"}
}

func MakeDailySections(list []UserEntry) []Section {
	sort.Stable(ByDate(list))

	dates := []string{}
	entries := map[string][]Entry{}
	for _, ue := range list {
		date := "No Diary"
		if !ue.mtime.IsZero() {
			date = ue.mtime.Format("2006/01/02")
		}

		list, ok := entries[date]
		if !ok {
			list = []Entry{}
			dates = append(dates, date)
		}
		entries[date] = append(list, NewDailyEntry(ue))
	}

	sections := []Section{}
	for _, date := range dates {
		sections = append(sections, Section{date, entries[date]})
	}

	return sections
}

func ListFiles(users []User) []UserEntry {
	list := []UserEntry{}
	for _, user := range users {
		ue := UserEntry{user: user}

		file := user.Path("", "diary.html")
		fi, err := os.Stat(path.Join(config.PrivateDir, file))
		if err == nil && fi.Mode().IsRegular() {
			ue.path = file
			ue.mtime = fi.ModTime()
		}

		list = append(list, ue)
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
