package main

import (
	"io"
	"os/exec"
)

type request struct {
	path    string
	w       io.Writer
	errChan chan error
}

type mkisofs struct {
	requestChan chan request
}

func initMkisofs() *mkisofs {
	mkiso := new(mkisofs)
	mkiso.requestChan = make(chan request, 10)

	worker := func(m *mkisofs) {
		for {
			select {
			case data := <-m.requestChan:
				cmd := exec.Command("mkisofs", "-R", "-V", "config-2", data.path)
				cmd.Stdout = data.w
				err := cmd.Run()
				if err != nil {
					data.errChan <- err
				} else {
					data.errChan <- nil
				}
			}
		}
	}

	//start some worker routines
	for i := 0; i < 5; i++ {
		go worker(mkiso)
	}

	return mkiso
}

func (m *mkisofs) create(path string, w io.Writer) error {
	err := make(chan error)
	req := request{
		path,
		w,
		err,
	}

	m.requestChan <- req
	return <-err
}
