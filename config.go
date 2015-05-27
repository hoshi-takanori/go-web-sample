package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/yosssi/ace"
)

type Config struct {
	Address    string
	PublicDir  string
	PrivateDir string

	Title string

	DatabaseDriver string
	DatabaseSource string

	SessionCookie http.Cookie

	AceOptions ace.Options
}

var config Config

func LoadConfig(filename string) error {
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(str, &config)
}
