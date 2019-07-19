package main

import (
	"io/ioutil"
	"net/http"
)

type down struct {
	baseHost string
}

func (d *down) BaseHost() string {
	return d.baseHost
}

func (d *down) Download(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
