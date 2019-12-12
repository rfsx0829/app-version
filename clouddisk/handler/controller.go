package handler

import (
	"encoding/json"
	"log"
	"net/http"

	mysql "github.com/rfsx0829/little-tools/clouddisk/mysql"
	util "github.com/rfsx0829/little-tools/clouddisk/util"
)

var (
	mp = map[string]string{
		"/":                 "help page",
		"/sign":             "sign up or log in",
		"/upload":           "upload file",
		"/user/{uid}":       "get user info by uid",
		"/file/{fid}":       "get file info by fid",
		"/files/{filehash}": "download a file by hash value",
	}
)

// Controller maintain handlers
type Controller struct {
	fileStorePath string
	db            mysql.Database
	tg            *util.TokenGenerater
}

// New create a Controller
func New(d mysql.Database, fileStorePath string) *Controller {
	return &Controller{fileStorePath, d, util.NewTokenGenerater()}
}

// HelpPageHandler make help page
func (c *Controller) HelpPageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(mp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("[help]:", err)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
