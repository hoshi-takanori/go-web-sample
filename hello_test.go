// +build hello

package main

import (
	"testing"
)

func TestMakeSections(t *testing.T) {
	println("TestListUsers")

	sections, err := MakeSections("nobody")
	if err != nil {
		panic(err)
	}

	for _, s := range sections {
		println(s.Name)
		for _, e := range s.Entries {
			println(e.Name, e.Path, e.Date, e.Dcls)
		}
		println()
	}
}
