package store

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"time"
	"os"

	"filestore/helper"
	"github.com/sirupsen/logrus"
)


type client struct {
	Logger     *logrus.Logger
	BaseURL string
	HttpClient *http.Client
}

// NewClient creates a new client
func NewClient(baseURL string) *client {
	return &client{
		BaseURL: baseURL,
		Logger: helper.NewLogger("filestore"),
		HttpClient: &http.Client {
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
			},
			Timeout: time.Second * 10,
		},
	}
}
// Add adds files to the store
func (c *client) Add(files []string) error {
	req, err := c.multipartBody(files)
	if err != nil {
		return err
	}
	resp, err := c.HttpClient.Do(req)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	return nil
}

// multipartBody adds files to a multipart writer
func (c *client) multipartBody(files []string) (*http.Request, error) {
	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)
	for _, fn := range files {
		file, err := os.Open(fn)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			return nil, err
		}
		part, err := bodyWriter.CreateFormFile(fn, fi.Name())
		if err != nil {
			return nil, err
		}
		io.Copy(part, file)
	}
	err := bodyWriter.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.BaseURL, bodyBuffer)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	return req, nil
}