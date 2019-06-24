package main

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type simpleHost struct {
	baseURL string
	token   string
}

func (s simpleHost) base() string {
	return s.baseURL
}

func (s simpleHost) up(file *os.File) error {
	// url := s.base() + "/upload"
	wter := multipart.NewWriter(nil)
	wter.CreateFormFile("file", file.Name())
	return nil
}

func (s simpleHost) down(file string) ([]byte, error) {
	url := s.base() + "/files/" + file
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}

func main() {
	sh := simpleHost{
		baseURL: "http://localhost:8000",
		token:   "token",
	}

	foo(sh)
}
