package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yosssi/ace"
)

type Config struct {
	Address string
	StaticDir string

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

	if config.StaticDir != "" {
		staticDir(config.StaticDir, http.DefaultServeMux)
	}

	http.HandleFunc("/", hello)
	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		panic(err)
	}
}

func staticDir(dirname string, mux *http.ServeMux) {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.Dir(dirname))
	for _, fi := range fis {
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			mux.Handle("/" + fi.Name() + "/", fileServer)
		}
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
