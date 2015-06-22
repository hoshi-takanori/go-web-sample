package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

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

func StaticDir(dirname string, mux *http.ServeMux) {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		println(err.Error())
		return
	}

	fileServer := http.FileServer(NoListFileSystem{http.Dir(dirname)})
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
