package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	bingURL    = "https://cn.bing.com"
	imgJSONURL = bingURL + "/HPImageArchive.aspx?format=js&idx=0&n=1"
)

func init() {
	folder := "/files"
	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		err = os.Mkdir(folder, os.ModePerm)
	}
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		saved = ""
		temp  []byte
	)

	routine := func() {
		today := time.Now().Format("2006_01_02")
		log.Println("[routine]", "today =", today)

		if today != saved {
			if data, err := getTodayImg(); err != nil {
				log.Println(err)
			} else {
				saved = today
				temp = data
				if err = ioutil.WriteFile("/files/"+saved+".jpg", temp, 0644); err != nil {
					log.Println(err)
				}
			}
		}
	}

	go func() {
		for {
			routine()
			time.Sleep(time.Hour * 6)
		}
	}()

	http.HandleFunc("/api/v1/bingimg", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[info]", r.Method, r.URL.Path)
		routine()
		w.WriteHeader(http.StatusOK)
		w.Write(temp)
	})

	port := "6666"
	if temp := os.Getenv("PORT"); temp != "" {
		port = temp
	}

	log.Println("Listeining on http://localhost:" + port + "/api/v1/bingimg ...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println(err)
	}
}

func getTodayImg() ([]byte, error) {
	url, err := getTodayImgURL()
	if err != nil {
		return nil, err
	}
	return getURLData(url)
}

func getTodayImgURL() (string, error) {
	data, err := getURLData(imgJSONURL)
	if err != nil {
		return "", err
	}

	var x struct {
		Images []map[string]interface{} `json:"images"`
	}
	// var x interface{}

	if err = json.Unmarshal(data, &x); err != nil {
		return "", err
	}

	imgURL := bingURL + x.Images[0]["url"].(string)
	// imgURL := x.(map[string]interface{})["images"].([]interface{})[0].(map[string]interface{})["url"].(string)
	return imgURL, nil
}

func getURLData(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(res.Body)
}
