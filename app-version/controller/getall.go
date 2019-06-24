package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/rfsx0829/little-tools/app-version/redis"
)

// GetAllProjects returns all projects
func (c *Controller) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	log.Println("[AP]", r.URL.Path)

	data, err := redis.GetSDataJSON(c.Projs)
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
func (c *Controller) GetAllVersions(w http.ResponseWriter, r *http.Request) {
	log.Println("[AV]", r.URL.Path)

	params := strings.Split(r.URL.Path, "/") // [ "", "{project}" ]
	project := params[1]

	data, err := redis.GetHDataJSON(project)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		log.Println("[AV]", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
