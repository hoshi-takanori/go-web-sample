package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/yosssi/ace"
)

type Config struct {
	Address string

	AceOptions ace.Options
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
	tpl, err := ace.Load("base", "hello", &config.AceOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, map[string]string{"Title": "Go Web Sample"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
