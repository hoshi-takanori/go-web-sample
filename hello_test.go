// +build hello

package main

import (
	"testing"
)

func TestMakeSections(t *testing.T) {
	println("TestListUsers")

	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	store, err = InitStore(config)
	if err != nil {
		panic(err)
	}

	sections, err := MakeSections(0)
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
