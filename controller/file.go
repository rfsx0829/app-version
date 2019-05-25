package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/rfsx0829/app-version/redis"
)

// GetFile read and return the file of a specified version
func (c Controller) GetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("[GF]", r.URL.Path)

	params := strings.Split(r.URL.Path, "/") // ["", "files", "{project}", "{version}"]
	project, version := params[2], params[3]

	data, err := c.getFile(project, version)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[GF]", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (c Controller) getFile(project, version string) ([]byte, error) {
	field := version + "_filepath"
	filePath, err := redis.Client.HGet(project, field).Result()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filePath)
}