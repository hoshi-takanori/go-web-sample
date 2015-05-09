package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Address string
}

var config Config

func main() {
	str, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(str, &config)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", hello)
	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, World!\n")
}
