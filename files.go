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

var editables = [...]string{".html", ".htm", ".css", ".js", ".java", ".go", ".md", ".txt"}

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
	dir := user.Path(config.PrivateDir, "")

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
	data["Dir"] = user.Path("", "")
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

	name := user.Path(config.PrivateDir, header.Filename)
	dst, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

	http.Redirect(w, r, "/files", http.StatusFound)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request, user User) {
	name := r.FormValue("name")
	if name == "" || strings.Contains(name, "/") {
		http.Error(w, "Bad Filename", http.StatusBadRequest)
		return
	}

	err := os.Truncate(user.Path(config.PrivateDir, name), 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/files", http.StatusFound)
}

func FileCopyHandler(w http.ResponseWriter, r *http.Request, user User) {
	src := r.FormValue("source")
	if src == "" || strings.Contains(src, "/") {
		http.Error(w, "Bad Filename", http.StatusBadRequest)
		return
	}

	srcPath := user.Path(config.PrivateDir, src)
	srcStat, err := os.Stat(srcPath)
	if err != nil || !srcStat.Mode().IsRegular() {
		http.Error(w, "Source Not Exists", http.StatusBadRequest)
		return
	}

	dst := r.FormValue("dest")
	if !GoodName(dst) {
		http.Error(w, "Bad Filename", http.StatusBadRequest)
		return
	}

	dstPath := user.Path(config.PrivateDir, dst)
	dstStat, err := os.Stat(dstPath)
	if !os.IsNotExist(err) &&
		(err != nil || !dstStat.Mode().IsRegular() || dstStat.Size() > 0) {
		http.Error(w, "Destination Exists", http.StatusBadRequest)
		return
	}

	reader, err := os.Open(srcPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	writer, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer writer.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/files", http.StatusFound)
}

func FileEditHandler(w http.ResponseWriter, r *http.Request, user User, name string) {
	if !Editable(name) {
		http.Error(w, "Bad Filename", http.StatusBadRequest)
		return
	}

	data := NewData("Edit " + name)
	data["Name"] = name
	ExecTemplate(w, "edit", data)
}

func FileGetHandler(w http.ResponseWriter, r *http.Request, user User, name string) {
	if !Editable(name) {
		http.Error(w, "Bad Filename", http.StatusBadRequest)
		return
	}

	file, err := os.Open(user.Path(config.PrivateDir, name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func FilePutHandler(w http.ResponseWriter, r *http.Request, user User, name string) {
	if !Editable(name) {
		http.Error(w, "Bad Filename", http.StatusBadRequest)
		return
	}

	dstPath := user.Path(config.PrivateDir, name)
	file, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
