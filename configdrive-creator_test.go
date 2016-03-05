package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

var b1 bytes.Buffer

func TestMkisofs(t *testing.T) {
	err := mkiso.create("testdata/bla", &b1)
	if err == nil {
		t.Error("err should not be nil")
	}
	b1.Reset()

	err = mkiso.create("testfiles/", &b1)
	if err != nil {
		t.Error(err)
	}
}

func TestApi(t *testing.T) {
	go main()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	data1 := b1.Bytes()

	filename := "testfiles/openstack/latest/user_data"
	url := "http://localhost:3000/configdrive"

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", "user_data")
	if err != nil {
		t.Error("error writing to buffer")
	}

	// open file handle
	fh, err := os.Open(filename)
	if err != nil {
		t.Error("error opening file")
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		t.Error(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if len(respBody) != len(data1) {
		t.Error("iso file has the wrong size")
	}

}
