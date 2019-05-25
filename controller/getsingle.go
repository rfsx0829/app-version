package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/rfsx0829/app-version/redis"
)

// Single handle the single version of a project
func (c Controller) Single(w http.ResponseWriter, r *http.Request) {
	var (
		data []byte
		err  error
	)
	log.Println("[SV]", r.URL.Path)

	params := strings.Split(r.URL.Path, "/") // [ "", "{project}", "{version}" ]
	project, version := params[1], params[2]

	switch {
	case !checkParam(params):
		err = errors.New("Invalid params")

	case strings.ToUpper(version) == "LATEST":
		data, err = c.getLatest(project)

	case r.Method == "POST":
		err = c.uploadFile(r, project, version)
		data = []byte("Upload Success !")

	default:
		data, err = c.getSingle(project, version)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
	log.Println("[SV]", "OK")
}

func (c Controller) getLatest(project string) ([]byte, error) {
	mp, err := redis.Client.HGetAll(project).Result()
	if err != nil {
		return nil, err
	}

	var x struct {
		Version string `json:"version"`
		URL     string `json:"url"`
	}

	for k, v := range mp {
		if later(k, x.Version) {
			x.Version = k
			x.URL = v
		}
	}

	return json.Marshal(x)
}

func (c Controller) getSingle(project, version string) ([]byte, error) {
	var x struct {
		Version string `json:"version"`
		URL     string `json:"url"`
	}
	log.Println("[GS]")

	res, err := redis.Client.HGet(project, version).Result()
	if err != nil {
		return nil, err
	}

	x.Version = version
	x.URL = res

	return json.Marshal(x)
}
