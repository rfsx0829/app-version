package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	mux "github.com/gorilla/mux"
	handlers "github.com/rfsx0829/little-tools/clouddisk/handler"
	mysql "github.com/rfsx0829/little-tools/clouddisk/mysql"
)

const (
	debug = false
)

const (
	baseHost string = "localhost"
	port     int    = 8999
)

func crossHandle(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func main() {
	db, err := mysql.Connect("localhost", 8806, "root", "123456")
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	if err = db.Initialize(); err != nil {
		log.Println(err)
		os.Exit(0)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		select {
		case <-ch:
			fmt.Println("Interrupt")
			if debug {
				db.DestroyTables()
			}
			os.Exit(0)
		}
	}()

	handle := handlers.New(db, "./files/")
	router := mux.NewRouter()

	router.HandleFunc("/", crossHandle(handle.HelpPageHandler))
	router.HandleFunc("/sign", crossHandle(handle.SignHandler))
	router.HandleFunc("/upload", crossHandle(handle.UploadHandler))
	router.HandleFunc("/user/{info}", crossHandle(handle.UserHandler))
	router.HandleFunc("/list", crossHandle(handle.FileListHandler))
	router.HandleFunc("/file/{args}", crossHandle(handle.FileHandler))
	router.HandleFunc("/files/{filehash}", crossHandle(handle.DownloadHandler))

	log.Printf("Serving on http://%s:%d", baseHost, port)
	log.Println(http.ListenAndServe(":"+strconv.Itoa(port), router))
}
