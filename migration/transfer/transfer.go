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
	Upload(name string, data []byte) error
}

// Downloader can download something
type Downloader interface {
	BaseURL() string
	Download(string) ([]byte, error)
}

// Transfer schedule the job
type Transfer struct {
	MainText []byte
	RegExp   string

	Uploader
	Downloader
}

func NewTrans() *Transfer {
	return &Transfer{}
}

func NewTransWithFilename(fileName string) (trans *Transfer, err error) {
	trans = NewTrans()

	trans.MainText, err = ioutil.ReadFile(fileName)
	return
}

func (t *Transfer) ChangeFile(newFileName string) (err error) {
	t.MainText, err = ioutil.ReadFile(newFileName)
	return
}

func (t *Transfer) Do() {
	links := t.FindAll()

	for _, e := range links {
		t.Once(string(e))
	}
}

func (t *Transfer) Once(link string) {
	fileName := link[strings.LastIndex(link, "/")+1:]

	if data, err := t.Download(link); err != nil {
		log.Println(err)
		log.Println(link)
	} else if err = t.Upload(fileName, data); err != nil {
		log.Println(err)
		log.Println(fileName, "failed. Write in local file.")
		if err = ioutil.WriteFile(fileName, data, os.ModePerm); err != nil {
			log.Println(err)
			log.Println("WriteFile failed.")
		}
	}
}

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
