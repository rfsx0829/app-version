package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rfsx0829/app-version/controller"
	"github.com/rfsx0829/app-version/redis"
)

type handler func(http.ResponseWriter, *http.Request)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

const (
	page404 string = `<html><head><title>哎呀，没有找到页面嗷</title></head><body>你来到了没有东西的荒原<br><br><a href="http://%s:%d">点击这里回首页嗷</a></body></html>`
)

func main() {
	const (
		redisHost string = "172.17.0.5"
		redisPort int    = 6379
		redisPass string = "pass"
		redisDB   int    = 0

		baseHost string = "localhost"
		port     int    = 8000
		projs    string = "_projs"
		pass     string = "token"
		rootDIR  string = "./"
	)

	redis.InitClient(redisHost, redisPort, redisPass, redisDB)

	router := mux.NewRouter()
	con := controller.New(baseHost, projs, pass, rootDIR, port)

	router.NotFoundHandler = handler(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(page404, baseHost, port)))
	})

	router.HandleFunc("/", con.GetAllProjects)
	router.HandleFunc("/{project}", con.GetAllVersions)
	router.HandleFunc("/{project}/{version}", con.Single)
	router.HandleFunc("/files/{project}/{version}", con.GetFile)

	log.Printf("Serving on http://%s:%d", baseHost, port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
