// +build entry

package main

import (
	"testing"
)

func TestMakeSections(t *testing.T) {
	println("TestMakeSections")

	year := 0
	order := true

	users, err := pgStore.ListUsers(year, order)
	if err != nil {
		panic(err)
	}

	list := ListFiles(users)
	var sections []Section
	if order {
		sections = MakeSections(list, year)
	} else {
		sections = MakeDailySections(list)
	}

	for _, s := range sections {
		println(s.Name)
		for _, e := range s.Entries {
			println(e.Name, e.Path, e.Date, e.Dcls)
		}
		println()
	}
}
