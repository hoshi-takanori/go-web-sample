package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	println("TestMain")

	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	store, err = InitStore(config)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	os.Exit(code)
}
