// +build postgres

package main

import (
	"testing"
)

func TestListUsers(t *testing.T) {
	println("TestListUsers")

	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	store, err = InitStore(config)
	if err != nil {
		panic(err)
	}

	users, err := pgStore.ListUsers(0)
	if err != nil {
		panic(err)
	}

	for _, u := range users {
		println(u.name, u.year, u.yearNo, u.staffYear)
	}
}
