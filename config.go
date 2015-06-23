package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/yosssi/ace"
)

type Config struct {
	Address    string
	PublicDir  string
	PrivateDir string

	Title string

	FreshYear  int
	FreshUntil string

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

	err = json.Unmarshal(str, &config)
	if err == nil {
		t, err := time.Parse("2006-01-02", config.FreshUntil)
		if err != nil || time.Now().After(t) {
			config.FreshYear = 0
		}
	}

	return err
}
