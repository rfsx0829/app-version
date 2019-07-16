package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type handler func(http.ResponseWriter, *http.Request)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

type controller struct {
	host  string
	port  int
	token string
	root  string
}

func (c *controller) Upload(w http.ResponseWriter, r *http.Request) {
	if fileName, err := c.upload(r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println(err)
	} else {
		retURL := fmt.Sprintf("http://%s:%d/files/%s", c.host, c.port, fileName)
		log.Println("[ret]", retURL)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(retURL))
	}
}

func (c *controller) upload(r *http.Request) (string, error) {
	if r.FormValue("token") != c.token {
		return "", errors.New("Invalid Token")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	suffix := strings.LastIndex(header.Filename, ".")
	fileName := hex.EncodeToString(hash(data)) + header.Filename[suffix:]

	log.Println("[upload]", fileName)

	if f, err := os.OpenFile(c.root+"files/"+fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm); err != nil {
		return "", err
	} else if _, err = io.Copy(f, bytes.NewBuffer(data)); err != nil {
		return "", err
	}

	return fileName, nil
}

func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	log.Println("[get]", c.root+r.URL.Path[1:])

	if data, err := ioutil.ReadFile(c.root + r.URL.Path[1:]); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println(err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func main() {
	const (
		baseHost string = "localhost"
		port     int    = 8000
		pass     string = "token"
		rootDIR  string = "./"
		page404  string = `<html><head><title>哎呀，没有找到页面嗷</title></head><body>你来到了没有东西的荒原<br><br><a href="http://%s:%d">点击这里回首页嗷</a></body></html>`
	)

	router := mux.NewRouter()
	con := controller{baseHost, port, pass, rootDIR}

	if err := createDir(rootDIR + "files"); err != nil {
		panic(err)
	}

	router.NotFoundHandler = handler(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(page404, baseHost, port)))
	})

	router.HandleFunc("/upload", con.Upload)
	router.HandleFunc("/files/{filename}", con.get)

	log.Printf("Serving on http://%s:%d", baseHost, port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func createDir(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		return err
	}

	return nil
}

func hash(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}
