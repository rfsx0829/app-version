package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rfsx0829/app-version/redis"
)

// Controller hold methods
type Controller struct{}

const (
	projs       string = "_projs"
	uploadToken string = "password"
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

	data, err := ioutil.ReadFile(r.URL.Path[:index] + "_" + r.URL.Path[index+1:])
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
	var (
		data []byte
		err  error
	)
	log.Println("[SV]", r.URL.Path)

	params := strings.Split(r.URL.Path, "/")
	if !checkParam(params) {
		err = errors.New("Invalid params")
	} else if strings.ToUpper(params[2]) == "LATEST" {
		data, err = c.getLatest(w, r, params[1])
	} else if r.Method == "POST" {
		err = c.uploadFile(r, params)
	} else {
		data, err = c.getSingle(params)
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

func (c Controller) getLatest(w http.ResponseWriter, r *http.Request, project string) ([]byte, error) {
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

func (c Controller) uploadFile(r *http.Request, params []string) error {
	if r.FormValue("token") != uploadToken {
		return errors.New("Invalid Token")
	}

	fileName := fmt.Sprintf("/files/%s_%s", params[1], params[2])
	url := fmt.Sprintf("http://39.98.162.91:8000/files/%s/%s", params[1], params[2])
	log.Println("[UF]", fileName)

	file, _, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
	}()

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	if _, err = io.Copy(f, file); err != nil {
		return err
	}

	redis.Client.SAdd(projs, params[1])
	if _, err := redis.Client.HSet(params[1], params[2], url).Result(); err != nil {
		return err
	}

	return nil
}

func (c Controller) getSingle(params []string) ([]byte, error) {
	var x struct {
		Version string `json:"version"`
		URL     string `json:"url"`
	}
	log.Println("[GS]")

	res, err := redis.Client.HGet(params[1], params[2]).Result()
	if err != nil {
		return nil, err
	}

	x.Version = params[2]
	x.URL = res

	return json.Marshal(x)
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

func later(newer, older string) bool {
	if len(older) == 0 {
		return true
	}

	if len(newer) == 0 {
		return false
	}

	var (
		i1     = strings.LastIndex(newer, "v")
		i2     = strings.LastIndex(older, "v")
		s1     = strings.Split(newer[i1+1:], ".")
		s2     = strings.Split(older[i2+1:], ".")
		length = min(len(s1), len(s2))
	)

	for i := 0; i < length; i++ {
		var (
			n1, _ = strconv.Atoi(s1[i])
			n2, _ = strconv.Atoi(s2[i])
		)

		if n1 == n2 {
			continue
		}

		if n1 > n2 {
			return true
		}
		return false
	}

	return false
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
