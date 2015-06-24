// +build user

package main

import (
	"testing"
)

func TestListUsers(t *testing.T) {
	println("TestListUsers")

	users, err := pgStore.ListUsers(0, true)
	if err != nil {
		panic(err)
	}

	for _, u := range users {
		println(u.name, u.year, u.yearNo, u.staffYear)
	}
}
