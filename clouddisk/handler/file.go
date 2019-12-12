package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	errs "github.com/rfsx0829/little-tools/clouddisk/errs"
	model "github.com/rfsx0829/little-tools/clouddisk/model"
	util "github.com/rfsx0829/little-tools/clouddisk/util"
)

type resInfo struct {
	StatusCode int         `json:"statuscode"`
	Data       interface{} `json:"data"`
}

// FileHandler deals file thing
func (c *Controller) FileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[file]:", r.URL.String())
	w.WriteHeader(http.StatusOK)
	var resInfo resInfo

	errFunc := func(err error) {
		resInfo.StatusCode = http.StatusBadRequest
		resInfo.Data = err.Error()
		resInfo.sendJSON(w)
		log.Println("[file]:", err)
	}
	if !strings.HasPrefix(strings.ToLower(r.URL.Path), "/file/") {
		errFunc(errs.ErrPath)
		return
	}
	fid, err := strconv.Atoi(r.URL.Path[6:])
	if err != nil {
		errFunc(err)
		return
	}
	results, err := c.db.SelectFile("*", "where fid="+strconv.Itoa(fid))
	if err != nil {
		errFunc(err)
		return
	}
	if len(results) != 1 {
		errFunc(errs.ErrWrongFID)
		return
	}
	f := results[0]
	f.Filepath = "******"
	resInfo.StatusCode = http.StatusOK
	resInfo.Data = f
	resInfo.sendJSON(w)
}

// DownloadHandler deals download thing
func (c *Controller) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[download]:", r.URL.String())
	w.WriteHeader(http.StatusOK)
	var resInfo resInfo

	errFunc := func(err error) {
		resInfo.StatusCode = http.StatusBadRequest
		resInfo.Data = err.Error()
		resInfo.sendJSON(w)
		log.Println("[download]:", err)
	}
	if !strings.HasPrefix(strings.ToLower(r.URL.Path), "/files/") {
		errFunc(errs.ErrPath)
		return
	}
	hashValue := r.URL.Path[7:]
	results, err := c.db.SelectFile("*", "where md5value='"+hashValue+"'")
	if err != nil {
		errFunc(err)
		return
	}
	if len(results) == 0 {
		errFunc(errs.ErrNoSuchFile)
		return
	}
	data, err := ioutil.ReadFile(results[0].Filepath)
	if err != nil {
		errFunc(err)
		return
	}
	w.Write(data)
}

// UploadHandler receive file from client
func (c *Controller) UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[upload]:", r.URL.String())
	w.WriteHeader(http.StatusOK)
	var resInfo resInfo

	errFunc := func(err error) {
		resInfo.StatusCode = http.StatusBadRequest
		resInfo.Data = err.Error()
		resInfo.sendJSON(w)
		log.Println("[upload]:", err)
	}

	if r.Method != "POST" {
		errFunc(errs.ErrPostOnly)
		return
	}

	token := r.FormValue("token")
	if token == "" {
		errFunc(errs.ErrTokenRequired)
		return
	}

	uid, err := c.tg.Token2UID(token)
	if err != nil {
		log.Println("[upload]:", err)
		errFunc(errs.ErrInvalidToken)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		errFunc(err)
		log.Println("???")
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		errFunc(err)
		return
	}
	hashValue := util.HashValue(data)

	f, err := os.OpenFile(c.fileStorePath+hashValue, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		errFunc(err)
		return
	}
	if _, err = f.Write(data); err != nil {
		errFunc(err)
		return
	}

	if err = c.db.InsertFile(&model.File{
		UID:      uid,
		Filename: header.Filename,
		Filepath: c.fileStorePath + hashValue,
		MD5Value: hashValue,
	}); err != nil {
		errFunc(err)
		return
	}

	resInfo.StatusCode = http.StatusOK
	resInfo.Data = "Success"
	resInfo.sendJSON(w)
}

// FileListHandler get file list
func (c *Controller) FileListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[list]:", r.URL.String())
	w.WriteHeader(http.StatusOK)
	var resInfo resInfo
	r.ParseForm()

	errFunc := func(err error) {
		resInfo.StatusCode = http.StatusBadRequest
		resInfo.Data = err.Error()
		resInfo.sendJSON(w)
		log.Println("[list]:", err)
	}

	uid, _ := strconv.Atoi(r.FormValue("uid"))
	token := r.FormValue("token")
	res, err := c.tg.Token2UID(token)
	if err != nil {
		errFunc(err)
		return
	}
	if res != uid {
		errFunc(errs.ErrInvalidToken)
		return
	}
	results, err := c.db.SelectFile("*", "where uid="+strconv.Itoa(uid))
	if err != nil {
		errFunc(err)
		return
	}

	for _, e := range results {
		e.Filepath = "******"
	}

	resInfo.StatusCode = http.StatusOK
	resInfo.Data = results
	resInfo.sendJSON(w)
}

func (r resInfo) sendJSON(w http.ResponseWriter) {
	data, _ := json.Marshal(r)
	w.Write(data)
}
