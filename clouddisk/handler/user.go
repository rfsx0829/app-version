package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	errs "github.com/rfsx0829/little-tools/clouddisk/errs"
	model "github.com/rfsx0829/little-tools/clouddisk/model"
	util "github.com/rfsx0829/little-tools/clouddisk/util"
)

type signInfo struct {
	IsSignUp bool   `json:"is_sign_up"`
	Email    string `json:"email"`
	Uname    string `json:"uname"`
	Upass    string `json:"upass"`
}

// SignHandler deals sign thing
func (c *Controller) SignHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[sign]:", r.URL.String())
	w.WriteHeader(http.StatusOK)
	var resInfo resInfo

	errFunc := func(err error) {
		resInfo.StatusCode = http.StatusBadRequest
		resInfo.Data = err.Error()
		resInfo.sendJSON(w)
		log.Println("[sign]:", err)
	}

	var x = &signInfo{}
	if err := jsonBody(r, x); err != nil {
		errFunc(err)
		return
	}

	if x.Uname == "" || x.Upass == "" {
		errFunc(errs.ErrNameRequired)
		return
	}

	str := fmt.Sprintf("where uname='%s'", x.Uname)
	results, err := c.db.SelectUser("*", str)
	if err != nil {
		errFunc(err)
		return
	}

	if x.IsSignUp {
		if len(results) != 0 {
			errFunc(errs.ErrNameUsed)
			return
		}
		if err := c.db.InsertUser(&model.User{
			Email:    x.Email,
			Uname:    x.Uname,
			Password: util.HashValue([]byte(x.Upass)),
		}); err != nil {
			errFunc(err)
			return
		}
		resInfo.Data = "SignUp OK"
	} else {
		if len(results) != 1 || results[0].Password != util.HashValue([]byte(x.Upass)) {
			errFunc(errs.ErrWrongNameOrPassword)
			return
		}
		token := c.tg.UID2Token(results[0].UID)
		resInfo.Data = struct {
			UID   int    `json:"uid"`
			Token string `json:"token"`
		}{results[0].UID, token}
	}

	resInfo.StatusCode = http.StatusOK
	resInfo.sendJSON(w)
}

// UserHandler get user info
func (c *Controller) UserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[user]:", r.URL.String())
	w.WriteHeader(http.StatusOK)
	var resInfo resInfo

	errFunc := func(err error) {
		resInfo.StatusCode = http.StatusBadRequest
		resInfo.Data = err.Error()
		resInfo.sendJSON(w)
		log.Println("[user]:", err)
	}
	if !strings.HasPrefix(strings.ToLower(r.URL.Path), "/user/") {
		errFunc(errs.ErrPath)
		return
	}
	uid, err := strconv.Atoi(r.URL.Path[6:])
	if err != nil {
		errFunc(err)
		return
	}
	results, err := c.db.SelectUser("*", "where uid="+strconv.Itoa(uid))
	if err != nil {
		errFunc(err)
		return
	}
	if len(results) != 1 {
		errFunc(errs.ErrWrongUID)
		return
	}
	u := results[0]
	u.Password = "******"
	resInfo.StatusCode = http.StatusOK
	resInfo.Data = u
	resInfo.sendJSON(w)
}

func jsonBody(r *http.Request, v interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
