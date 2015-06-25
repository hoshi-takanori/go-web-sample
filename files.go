package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type FileEntry struct {
	Name string
	Edit bool
	Size int64
	Date string
}

func FilesHandler(w http.ResponseWriter, r *http.Request, user User) {
	dir := user.Dir(config.PrivateDir, "")

	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files := []FileEntry{}
	for _, fi := range fis {
		if fi.Mode().IsRegular() && !strings.HasPrefix(fi.Name(), ".") {
			edit := true
			date := fi.ModTime().Format("2006/01/02 15:04:05")
			files = append(files, FileEntry{fi.Name(), edit, fi.Size(), date})
		}
	}

	data := NewData("File Manager")
	data["Username"] = user.name
	data["Dir"] = user.Dir("", "")
	data["Files"] = files
	ExecTemplate(w, "files", data)
}

func FileUploadHandler(w http.ResponseWriter, r *http.Request, user User) {
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request, user User) {
}

func FileCopyHandler(w http.ResponseWriter, r *http.Request, user User) {
}

func FileEditHandler(w http.ResponseWriter, r *http.Request, user User, name string) {
}

func FileGetHandler(w http.ResponseWriter, r *http.Request, user User, name string) {
}

func FilePutHandler(w http.ResponseWriter, r *http.Request, user User, name string) {
}
