package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var editables = [...]string{".html", ".htm", ".css", ".js", ".java", ".txt"}

type FileEntry struct {
	Name string
	Edit bool
	Size int64
	Date string
}

func GoodName(name string) bool {
	match, err := regexp.MatchString("^\\w[-.\\w]*$", name)
	return err == nil && match
}

func Editable(name string) bool {
	if !GoodName(name) {
		return false
	}
	ext := path.Ext(name)
	for _, e := range editables {
		if e == ext {
			return true
		}
	}
	return false
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
			edit := Editable(fi.Name())
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
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header != nil && !GoodName(header.Filename) {
		http.Error(w, "Bad Filename", http.StatusBadRequest)
		return
	}

	size, err := file.Seek(0, 2)
	if err == nil {
		_, err = file.Seek(0, 0)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if size > 1048576 {
		http.Error(w, "File Too Large", http.StatusBadRequest)
		return
	}

	name := user.Dir(config.PrivateDir, header.Filename)
	dst, err := os.OpenFile(name, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/files", 302)
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
