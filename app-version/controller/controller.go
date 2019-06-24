package controller

import (
	"os"
	"strconv"
	"strings"
)

// Controller hold methods
type Controller struct {
	Host        string
	Port        int
	Projs       string
	UploadToken string
	Root        string
}

// New create a Controller
func New(h, p, t, r string, port int) *Controller {
	err := createDir(r + "files")
	if err != nil {
		panic(err)
	}

	return &Controller{
		Host:        h,
		Port:        port,
		Projs:       p,
		UploadToken: t,
		Root:        r,
	}
}

func createDir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
	}

	return err
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
