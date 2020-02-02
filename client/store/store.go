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

// newStoreRequest builds a request for the store server
func newStoreRequest(method, url string, body *bytes.Buffer) (*http.Request, error){
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}
	return req, nil
}

// multipartBody adds files to a multipart writer
func multipartBody(files []string) (*bytes.Buffer, string, error) {
	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)
	for _, fn := range files {
		file, err := os.Open(fn)
		if err != nil {
			return nil, "", err
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			return nil, "", err
		}
		part, err := bodyWriter.CreateFormFile(fn, fi.Name())
		if err != nil {
			return nil, "", err
		}
		io.Copy(part, file)
	}
	err := bodyWriter.Close()
	if err != nil {
		return nil, "", err
	}
	return bodyBuffer, bodyWriter.FormDataContentType(), nil
//	req, err := http.NewRequest("POST", fmt.Sprintf("%s/add",c.BaseURL), bodyBuffer)
//	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
//	return req, nil
}

// Add adds files to the store
func (c *client) Add(files []string) error {
//	req, err := c.multipartBody(files)
//	c.Logger.Infof("request %v", req)
//	if err != nil {
//		return err
//	}
	bodyBuffer, contentType, err := multipartBody(files) 
	if err != nil {
		c.Logger.Fatalf("Could not build request %v", err)
		return err
	}

	req, err := newStoreRequest("POST", fmt.Sprintf("%s/add", c.BaseURL), bodyBuffer)
	c.Logger.Debugf("req %v", req)
	if err != nil {
		c.Logger.Fatalf("Could not build request %v", err)
		return err
	}
	req.Header.Add("Content-Type", contentType)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) || (err != nil) {
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	return nil
}

// Remove removes the file from the store
func (c *client) Remove(file string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/remove?file=%s", c.BaseURL, file), nil)
	c.Logger.Debugf("req %v", req)
	if err != nil {
		c.Logger.Fatalf("Could not build request %v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	return nil
}

// List lists all files in the store
func (c *client) List() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/list", c.BaseURL), nil)
	c.Logger.Debugf("request %v", req)
	if err != nil {
		return err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	fmt.Fprintf(os.Stdout, string(b))
	return nil
}

// Update updates or create a file in the store
func (c *client) Update(file string) error {
	bodyBuffer, contentType, err := multipartBody([]string{file})
	if err != nil {
		c.Logger.Fatalf("Could not read body %v", err)
		return err
	}
	req, err := newStoreRequest("POST", fmt.Sprintf("%s/update", c.BaseURL), bodyBuffer)
	c.Logger.Debugf("request %v", req)
	if err != nil {
		c.Logger.Fatalf("Could not build request %v", err)
		return err
	}
	req.Header.Add("Content-Type", contentType)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	return nil
}