package controller

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/rfsx0829/little-tools/app-version/redis"
)

func (c *Controller) uploadFile(r *http.Request, project, version string) error {
	if r.FormValue("token") != c.UploadToken {
		return errors.New("Invalid Token")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
	}()

	fileName := fmt.Sprintf("%sfiles/%s_%s_%s", c.Root, project, version, header.Filename)
	retURL := fmt.Sprintf("http://%s:%d/files/%s/%s_%s", c.Host, c.Port, project, version, header.Filename)
	log.Println("[UF]", fileName)

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	if _, err = io.Copy(f, file); err != nil {
		return err
	}

	mp := map[string]interface{}{
		version:               retURL,
		version + "_filepath": fileName,
	}

	_, err = redis.Client.SAdd(c.Projs, project).Result()
	if err != nil {
		return err
	}

	_, err = redis.Client.HMSet(project, mp).Result()
	return err
}
