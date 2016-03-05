package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/xlog"
	"golang.org/x/net/context"
)

func indexHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("./static/index.html")
	w.Write(body)
}

func configdriveHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	l := xlog.FromContext(ctx)

	//The whole request body is parsed and up to a total of 34MB bytes of its file parts are stored in memory,
	//with the remainder stored on disk in temporary files.
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		l.Error(err)
		return
	}
	defer file.Close()

	//add the filename to access log
	l.SetField("filename", handler.Filename)
	tempDir, err := ioutil.TempDir(os.TempDir(), "mkisofs")
	if err != nil {
		l.Error(err)
		return
	}
	defer os.RemoveAll(tempDir)

	//Create the folder structure inside the temp folder
	path := filepath.Join(tempDir, "openstack", "latest")
	os.MkdirAll(path, os.ModePerm)
	f, err := os.OpenFile(filepath.Join(path, "user_data"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(f, file)
	f.Close()

	err = mkiso.create(tempDir, w)
	if err != nil {
		l.Error(err)
		return
	}
}
