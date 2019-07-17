package transfer

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// Uploader can upload something
type Uploader interface {
	Upload(name string, data []byte) (link string, err error)
}

// Downloader can download something
type Downloader interface {
	BaseURL() string
	Download(string) ([]byte, error)
}

// Transfer schedule the job
type Transfer struct {
	FileName string
	MainText []byte
	RegExp   string

	Uploader
	Downloader
}

// NewTrans create a new Transfer, and return a pointer which point to the Transfer
func NewTrans() *Transfer {
	return &Transfer{}
}

// NewTransWithFilename create a new Transfer and set its MainText to localfile's data.
func NewTransWithFilename(fileName string) (*Transfer, error) {
	trans := NewTrans()

	return trans, trans.ChangeFile(fileName)
}

// ChangeFile change MainText by localfile's data
func (t *Transfer) ChangeFile(newFileName string) (err error) {
	t.MainText, err = ioutil.ReadFile(newFileName)
	if err == nil {
		t.FileName = newFileName
	}
	return
}

// Once just deal one link
func (t *Transfer) Once(link string) (string, error) {
	fileName := link[strings.LastIndex(link, "/")+1:]

	data, err := t.Download(link)
	if err != nil {
		return "", err
	}

	newLink, err := t.Upload(fileName, data)
	if err != nil {
		log.Println(fileName, "failed. Try to write in local file.")
		if err2 := ioutil.WriteFile(fileName, data, os.ModePerm); err2 != nil {
			log.Println("Write file failed.", err2)
			return "", err
		}
	}

	return newLink, nil
}

// FindAll find all links that match r.RegExp and contains baseURL
func (t *Transfer) FindAll() [][]byte {
	baseURL := t.BaseURL()
	results := regexp.MustCompile(t.RegExp).FindAll(t.MainText, -1)

	for i, e := range results {
		if !strings.Contains(string(e), baseURL) {
			results = append(results[:i], results[i+1:]...)
		}
	}

	return results
}
