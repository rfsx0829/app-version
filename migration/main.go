package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/rfsx0829/little-tools/migration/transfer"
)

type bed struct {
	baseHost string
	token    string
}

func (b *bed) Upload(fileName string, data []byte) (string, error) {
	buf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(buf)
	bodyWriter.SetBoundary("qweadqweas")

	bodyWriter.WriteField("token", b.token)

	if _, err := bodyWriter.CreateFormFile("file", fileName); err != nil {
		return "", err
	}
	buf.Write(data)

	bodyWriter.Close()

	req, err := http.NewRequest("POST", b.baseHost, io.MultiReader(buf))
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

type pair struct {
	old string
	new string
}

func main() {
	trans, err := transfer.NewTransWithFilename("./xxx.md")
	if err != nil {
		panic(err)
	}

	trans.RegExp = "http[s?]://[\\d\\w\\./-]*\\.(jpg|png)"
	trans.Uploader = &bed{"http://localhost:8000/upload", "token"}
	trans.Downloader = &down{"jianshu.io"}

	oldLinks := trans.FindAll()
	pairs := make([]*pair, 0, len(oldLinks))

	for _, e := range oldLinks {
		pairs = append(pairs, &pair{
			old: string(e),
			new: "",
		})
	}

	for _, e := range pairs {
		if newLink, err := trans.Once(e.old); err != nil {
			log.Println(err)
		} else {
			e.new = newLink
		}
	}

	trans.MainText = replace(trans.MainText, pairs)

	index := strings.LastIndex(trans.FileName, ".")
	newFileName := trans.FileName[:index] + "_new" + trans.FileName[index:]

	if err = ioutil.WriteFile(newFileName, trans.MainText, os.ModePerm); err != nil {
		panic(err)
	}
}

func replace(text []byte, pairs []*pair) []byte {
	newStr := string(text)

	for _, e := range pairs {
		if e.new != "" {
			newStr = strings.ReplaceAll(newStr, e.old, e.new)
		}
	}

	return []byte(newStr)
}
