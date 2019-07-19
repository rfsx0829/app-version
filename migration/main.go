package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/rfsx0829/little-tools/migration/transfer"
)

func main() {
	var (
		configFile = flag.String("cfg", "", "[optional] config file, if this specified, other flags won't be used.")
		inputFile  = flag.String("if", "", "[*] input file")
		outFile    = flag.String("of", "", "[optional] output file")
		regExp     = flag.String("re", "http[s?]://[\\d\\w\\./-]*\\.(jpg|png)", "[optional] regular expression used to match link.")
		tokenKey   = flag.String("tk", "token", "[optional] tokenKey, use default if you dont know what it mean.")
		fileKey    = flag.String("fk", "file", "[optional] fileKey, use default if you dont know what it mean.")
		token      = flag.String("token", "token", "[optional] the token if server needed")
		uploadURL  = flag.String("up", "", "[*] the url that upload picture. ex: http://localhost:8000/upload")
		baseURL    = flag.String("down", "", "[*] the base url that you download picture. ex: jianshu.io")

		trans = transfer.NewTrans()
	)

	flag.CommandLine.Usage = flag.PrintDefaults
	flag.Parse()

	if *configFile != "" {
		// TODO: json.Unmarshal can't deal interface. I guess.
		if data, err := ioutil.ReadFile(*configFile); err != nil {
			fmt.Println(err)
			return
		} else if err = json.Unmarshal(data, trans); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if *uploadURL == "" {
			fmt.Println("No upload URL.")
			return
		}
		if *baseURL == "" {
			fmt.Println("No base URL.")
			return
		}
		if *inputFile == "" {
			fmt.Println("No input file.")
			return
		}

		trans.OutFile = *outFile
		trans.RegExp = *regExp
		trans.Uploader = &bed{
			uploadURL: *uploadURL,
			token:     *token,
			tokenKey:  *tokenKey,
			fileKey:   *fileKey,
		}
		trans.Downloader = &down{
			baseHost: *baseURL,
		}
	}

	info, err := os.Stat(*inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if info.IsDir() {
		if fd, err := ioutil.ReadDir(*inputFile); err == nil {
			for _, e := range fd {
				if !strings.HasSuffix(e.Name(), ".md") {
					continue
				}
				if err = changeAndRun(trans, path.Join(*inputFile, e.Name())); err != nil {
					fmt.Println(err)
					return
				}
			}
		} else {
			fmt.Println(err)
			return
		}
	} else {
		if err = changeAndRun(trans, *inputFile); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func changeAndRun(t *transfer.Transfer, fileName string) (err error) {
	if err = t.ChangeFile(fileName, ""); err != nil {
		return err
	}
	return t.Run()
}
