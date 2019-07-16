package transfer

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// Uploader can upload something
type Uploader interface {
	Upload([]byte) error
}

// Downloader can download something
type Downloader interface {
	BaseURL() string
	Download(string) ([]byte, error)
}

// Transfer schedule the job
type Transfer struct {
	MainText []byte

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

func (t *Transfer) Do() {
	links := t.findAll()
	for _, e := range links {
		if data, err := t.Download(string(e)); err != nil {
			log.Println(err)
		} else if err = t.Upload(data); err != nil {
			log.Println(err)
		}
	}
}

func (t *Transfer) findAll() [][]byte {
	baseURL := t.BaseURL()
	results := regexp.MustCompile("https?://\\S+\\.(jpg|png)").FindAll(t.MainText, -1)
	for i, e := range results {
		if !strings.Contains(string(e), baseURL) {
			results = append(results[:i], results[i+1:]...)
		}
	}

	return results
}
