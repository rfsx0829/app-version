package main

import (
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

func main() {
	redis.InitClient("192.168.99.105", 31150, "password", 0)

	router := mux.NewRouter()
	con := controller.Controller{}

	router.NotFoundHandler = handler(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello 404"))
	})
	router.HandleFunc("/", con.GetAllProjects)
	router.HandleFunc("/{project}", con.GetAllVersions)
	router.HandleFunc("/{project}/{version}", con.Single)
	router.HandleFunc("/files/{project}/{version}", con.GetFile)

	log.Println("Serving on localhost:8000")
	log.Println(http.ListenAndServe(":8000", router))
}
