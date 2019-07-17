package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/rfsx0829/little-tools/migration/transfer"
)

type bed struct {
	baseHost string
	token    string
}

func (b *bed) Upload(fileName string, data []byte) (err error) {
	buf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(buf)

	if err = bodyWriter.SetBoundary("qweadqweas"); err != nil {
		log.Println(err)
		return
	}

	if err = bodyWriter.WriteField("token", b.token); err != nil {
		log.Println(err)
		return
	}

	if _, err = bodyWriter.CreateFormFile("file", fileName); err != nil {
		log.Println(err)
		return
	}

	if _, err = buf.Write(data); err != nil {
		log.Println(err)
		return
	}

	bodyWriter.Close()

	reqReader := io.MultiReader(buf)
	req, err := http.NewRequest("POST", b.baseHost, reqReader)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Connection", "close")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	req.ContentLength = int64(buf.Len())

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if data, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(data))
	}

	return nil
}

type down struct {
	baseHost string
}

func (d *down) BaseURL() string {
	return d.baseHost
}

func (d *down) Download(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func main() {
	trans, err := transfer.NewTransWithFilename("./xxx.md")
	if err != nil {
		panic(err)
	}

	trans.RegExp = "http[s?]://[\\d\\w\\./-]*\\.(jpg|png)"
	trans.Uploader = &bed{"http://localhost:8000/upload", "token"}
	trans.Downloader = &down{"jianshu.io"}

	trans.Do()
}
