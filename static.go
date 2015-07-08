package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/russross/blackfriday"
)

var CustomTypeMap = map[string]string{
	"text/x-java": "text/plain",
}

type CustomResponseWriter struct {
	w http.ResponseWriter
}

func (cw CustomResponseWriter) Header() http.Header {
	return cw.w.Header()
}

func (cw CustomResponseWriter) Write(b []byte) (int, error) {
	return cw.w.Write(b)
}

func (cw CustomResponseWriter) WriteHeader(code int) {
	contentType := cw.Header().Get("Content-Type")
	for k, v := range CustomTypeMap {
		if contentType == k || strings.HasPrefix(contentType, k+";") {
			cw.Header().Set("Content-Type", strings.Replace(contentType, k, v, 1))
			break
		}
	}

	cw.w.WriteHeader(code)
}

type NoListFile struct {
	http.File
}

func (f NoListFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

type NoListFileSystem struct {
	base http.FileSystem
}

func (fs NoListFileSystem) Open(name string) (http.File, error) {
	f, err := fs.base.Open(name)
	if err != nil {
		return nil, err
	}
	return NoListFile{f}, nil
}

func HandleMarkDown(w http.ResponseWriter, r *http.Request, dir http.Dir) bool {
	if !strings.HasSuffix(r.URL.Path, ".md") {
		return false
	}

	file, err := dir.Open(r.URL.Path)
	if err != nil {
		return false
	}
	defer file.Close()

	md, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}

	w.Write(blackfriday.MarkdownCommon(md))
	return true
}

func CustomFileServer(dir http.Dir) http.HandlerFunc {
	fileServer := http.FileServer(NoListFileSystem{dir})
	return func(w http.ResponseWriter, r *http.Request) {
		if !HandleMarkDown(w, r, dir) {
			fileServer.ServeHTTP(CustomResponseWriter{w}, r)
		}
	}
}

func StaticDir(dirname string, mux *http.ServeMux) {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		println(err.Error())
		return
	}

	fileServer := CustomFileServer(http.Dir(dirname))
	for _, fi := range fis {
		if fi.Mode()&os.ModeSymlink != 0 {
			fi2, err := os.Stat(path.Join(dirname, fi.Name()))
			if err == nil {
				fi = fi2
			}
		}
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			mux.Handle("/"+fi.Name()+"/", fileServer)
		}
	}
}
