package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rfsx0829/app-version/redis"
)

// Controller hold methods
type Controller struct{}

const (
	projs string = "_projs"
)

// GetAllProjects returns all projects
func (c Controller) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	log.Println("[AP]", r.URL.Path)

	data, err := redis.GetSDataJSON(projs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[AP]", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetAllVersions of a project
func (c Controller) GetAllVersions(w http.ResponseWriter, r *http.Request) {
	log.Println("[AV]", r.URL.Path)

	projectName := r.URL.Path[1:]

	data, err := redis.GetHDataJSON(projectName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[AV]", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetFile ret file
func (c Controller) GetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("[GF]", r.URL.Path)

	index := strings.LastIndex(r.URL.Path, "/")

	data, err := ioutil.ReadFile("." + r.URL.Path[:index] + "_" + r.URL.Path[index+1:])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[GF]", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Single handle the single version of a project
func (c Controller) Single(w http.ResponseWriter, r *http.Request) {
	log.Println("[SV]", r.URL.Path)
	params := strings.Split(r.URL.Path, "/")
	if !checkParam(params) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Params"))
		log.Println("[SV]", "Invalid Params")
		return
	}

	if r.Method == "POST" {
		c.uploadFile(w, r, params)
		return
	}

	c.getSingle(w, r, params)
}

func (c Controller) uploadFile(w http.ResponseWriter, r *http.Request, params []string) {
	// TODO: /files
	fileName := fmt.Sprintf("./files/%s_%s", params[1], params[2])
	url := fmt.Sprintf("http://localhost:8000/files/%s/%s", params[1], params[2])
	log.Println("[UF]", fileName)

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[UF]", err)
		return
	}
	defer func() {
		file.Close()
	}()

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[UF]", err)
		return
	}

	if _, err = io.Copy(f, file); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[UF]", err)
		return
	}

	redis.Client.SAdd(projs, params[1])
	if _, err := redis.Client.HSet(params[1], params[2], url).Result(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[UF]", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upload Success !"))
	log.Println("[UF]", "Success !")
}

func (c Controller) getSingle(w http.ResponseWriter, r *http.Request, params []string) {
	var x struct {
		Version string `json:"version"`
		URL     string `json:"url"`
	}
	log.Println("[GS]")

	res, err := redis.Client.HGet(params[1], params[2]).Result()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[GS]", err)
		return
	}

	x.Version = params[2]
	x.URL = res

	data, err := json.Marshal(x)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[GS]", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func checkParam(params []string) bool {
	if len(params) < 3 {
		return false
	}

	if len(params[1]) < 1 {
		return false
	}

	if len(params[2]) < 1 {
		return false
	}

	return true
}
