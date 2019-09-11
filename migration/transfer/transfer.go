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
	BaseHost() string
	Download(string) ([]byte, error)
}

type pair struct {
	old string
	new string
}

// Transfer schedule the job
type Transfer struct {
	FileName string
	OutFile  string
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

	return trans, trans.ChangeFile(fileName, "")
}

// ChangeFile change MainText by localfile's data
func (t *Transfer) ChangeFile(newFileName string, newOutFile string) (err error) {
	t.MainText, err = ioutil.ReadFile(newFileName)
	if err == nil {
		t.FileName = newFileName
	}
	t.OutFile = newOutFile
	return
}

// Once just deal one link
func (t *Transfer) Once(link string) (string, error) {
	fileName := link[strings.LastIndex(link, "/")+1:]

	log.Println("download ... ", link)
	data, err := t.Download(link)
	if err != nil {
		return "", err
	}

	log.Println("upload   ... ", fileName)
	newLink, err := t.Upload(fileName, data)
	if err != nil {
		log.Println(fileName, "failed. Try to write in local file.")
		_ = ioutil.WriteFile(fileName, data, os.ModePerm)
		return "", err
	}

	return newLink, nil
}

// FindLinks find all links that match r.RegExp and contains baseURL
func (t *Transfer) FindLinks() [][]byte {
	baseURL := t.BaseHost()
	results := regexp.MustCompile(t.RegExp).FindAll(t.MainText, -1)
	flitered := make([][]byte, 0, len(results))

	for _, e := range results {
		if strings.Contains(string(e), baseURL) {
			flitered = append(flitered, e)
		}
	}

	return flitered
}

// Run the transfer
func (t *Transfer) Run() error {
	oldLinks := t.FindLinks()
	if len(oldLinks) == 0 {
		log.Println("No matched pictures. --- ", t.FileName)
		return t.WriteFile()
	}
	pairs := make([]*pair, 0, len(oldLinks))

	for _, e := range oldLinks {
		pairs = append(pairs, &pair{
			old: string(e),
			new: "",
		})
	}

	for _, e := range pairs {
		if newLink, err := t.Once(e.old); err == nil {
			e.new = newLink
		} else {
			return err
		}
	}

	t.replace(pairs)

	return t.WriteFile()
}

// WriteFile write t.MainText to t.OutFile
func (t *Transfer) WriteFile() error {
	if t.OutFile == "" {
		index := strings.LastIndex(t.FileName, ".")
		t.OutFile = t.FileName[:index] + "_new" + t.FileName[index:]
	}

	return ioutil.WriteFile(t.OutFile, t.MainText, 0740)
}

func (t *Transfer) replace(pairs []*pair) {
	newStr := string(t.MainText)

	for _, e := range pairs {
		if e.new != "" {
			newStr = strings.ReplaceAll(newStr, e.old, e.new)
		}
	}

	t.MainText = []byte(newStr)
}
