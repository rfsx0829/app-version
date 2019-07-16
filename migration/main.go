package main

import (
	"fmt"

	"github.com/rfsx0829/little-tools/migration/transfer"
)

type bed struct {
	baseHost string
}

func (b *bed) Upload(data []byte) error {
	fmt.Println("upload", string(data))
	return nil
}

type down struct {
	baseHost string
}

func (d *down) BaseURL() string {
	return d.baseHost
}

func (d *down) Download(link string) ([]byte, error) {
	fmt.Println("download", link)
	return nil, nil
}

func main() {
	trans, err := transfer.NewTransWithFilename("./xxx.md")
	if err != nil {
		panic(err)
	}

	trans.Uploader = &bed{""}
	trans.Downloader = &down{"jianshu.io"}

	trans.Do()
}
