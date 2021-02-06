package http

import (
	"bytes"
	"fmt"
	"github.com/brickman-source/golang-utilities/log"
	"mime/multipart"
)

func getFirstFileFromMultipartForm(form *multipart.Form) ([]byte, error) {
	for _, files := range form.File {
		for _, file := range files {
			src, err := file.Open()
			if err != nil {
				log.Errorf("open file err: %s", err.Error())
				return nil, err
			}
			defer src.Close()

			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(src)
			if err != nil {
				log.Errorf("open file err: %s", err.Error())
				return nil, err
			}
			// 只取第一个文件
			return buf.Bytes(), nil
		}
	}
	return nil, fmt.Errorf("no file found")
}
