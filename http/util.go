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
			ret := buf.Bytes()
			log.Infof("file size: %v", len(ret))
			return ret, nil
		}
	}
	return nil, fmt.Errorf("no file found")
}

func getFileFromMultipartForm(form *multipart.Form, fieldName string) ([]byte, error) {
	for k, files := range form.File {
		if fieldName != k {
			continue
		}
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
			ret := buf.Bytes()
			log.Infof("file size: %v", len(ret))
			return ret, nil
		}
	}
	return nil, fmt.Errorf("no file found")
}

func getValueFromMultipartForm(form *multipart.Form, fieldName string) (string, error) {
	for k, values := range form.Value {
		if k == fieldName && len(values) > 0 {
			return values[0], nil
		}
	}
	return "", fmt.Errorf("no value found")
}

