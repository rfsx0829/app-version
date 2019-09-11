package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

type bed struct {
	uploadURL string
	token     string
	tokenKey  string
	fileKey   string
}

// Upload data to uploadURL
func (b *bed) Upload(fileName string, data []byte) (string, error) {
	if b.checkURL() {
		return b.upload(fileName, data)
	}

	if strings.Index(b.uploadURL, "./") != 0 || strings.Index(b.uploadURL, ".") != 0 {
		b.uploadURL = "./" + b.uploadURL
	}

	return "", ioutil.WriteFile(b.uploadURL+fileName, data, 0640)
}

func (b *bed) upload(fileName string, data []byte) (string, error) {
	buf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(buf)
	bodyWriter.SetBoundary("shoule be garbled ?")

	bodyWriter.WriteField(b.tokenKey, b.token)

	if _, err := bodyWriter.CreateFormFile(b.fileKey, fileName); err != nil {
		return "", err
	}
	buf.Write(data)

	bodyWriter.Close()

	req, err := http.NewRequest("POST", b.uploadURL, io.MultiReader(buf))
	if err != nil {
		return "", err
	}

	req.Header.Set("Connection", "close")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	req.ContentLength = int64(buf.Len())

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respData), nil
}

func (b *bed) checkURL() bool {
	reg := regexp.MustCompile("http[s]?://\\S+")
	return reg.MatchString(b.uploadURL)
}
