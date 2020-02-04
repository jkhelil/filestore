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
	"github.com/spf13/viper"
)

// Client interface for basic client structure
type Client struct {
	Logger     *logrus.Logger
	BaseURL string
	HTTPClient *http.Client
}

// NewClient creates a new client
func NewClient() *Client {
	return &Client{
		BaseURL: viper.GetString("server-url"),
		Logger: helper.NewLogger("filestore"),
		HTTPClient: &http.Client {
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
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, "", err
		}
	}
	err := bodyWriter.Close()
	if err != nil {
		return nil, "", err
	}
	return bodyBuffer, bodyWriter.FormDataContentType(), nil
}

// Add adds files to the store
func (c *Client) Add(files []string) error {
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

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) || (err != nil) {
		c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	c.Logger.Infof("HTTPStatusCode: '%d'; ResponseMessage: '%s'", resp.StatusCode, string(b))
	return nil
}

// Remove removes the file from the store
func (c *Client) Remove(file string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/remove?file=%s", c.BaseURL, file), nil)
	c.Logger.Debugf("req %v", req)
	if err != nil {
		c.Logger.Fatalf("Could not build request %v", err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	c.Logger.Infof("HTTPStatusCode: '%d'; ResponseMessage: '%s'", resp.StatusCode, string(b))
	return nil
}

// List lists all files in the store
func (c *Client) List() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/list", c.BaseURL), nil)
	c.Logger.Debugf("request %v", req)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	c.Logger.Debugf("HTTPStatusCode: '%d'", resp.StatusCode)
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}

// Update updates or create a file in the store
func (c *Client) Update(file string) error {
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

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	c.Logger.Infof("HTTPStatusCode: '%d'; ResponseMessage: '%s'", resp.StatusCode, string(b))
	return nil
}


// FreqWords prints most 10 frequent words
func (c *Client) FreqWords() error {
	limit := viper.GetInt("limit")
	order := viper.GetString("order")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/freqwords?limit=%d&order=%s", c.BaseURL, limit, order), nil)
	c.Logger.Debugf("request %v", req)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	c.Logger.Debugf("HTTPStatusCode: '%d'", resp.StatusCode)
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}

// CountWords prints most 10 frequent words
func (c *Client) CountWords() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/countwords", c.BaseURL), nil)
	c.Logger.Debugf("request %v", req)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Logger.Fatalf("Could not get response %v", err)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) || (err != nil) {
		c.Logger.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
		return fmt.Errorf("HTTPStatusCode: '%d'; ResponseMessage: '%s'; ErrorMessage: '%v'", resp.StatusCode, string(b), err)
	}
	c.Logger.Debugf("HTTPStatusCode: '%d'", resp.StatusCode)
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}