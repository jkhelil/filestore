package filestore

import (
	"io"
	"os"
	"testing"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
)

func TestAdd(t *testing.T) {
	config := NewConfig()
	config.StoreDir = "./testdata"
	fs := NewFileStore(config)
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()
		part, _ := writer.CreateFormFile("file", "test.txt")
		content := []byte("Test Add function")
		part.Write(content)
	}()

	req := httptest.NewRequest("POST", "/", pr)
    req.Header.Add("Content-Type", writer.FormDataContentType())
	handler := func(w http.ResponseWriter, r *http.Request) {
		fs.Add(w, r)
	}
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	t.Log("It should respond with an HTTP status code of 200")
    if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
        t.Errorf("Expected %d, received %d", 201, resp.StatusCode)
    }
    t.Log("It should create a file named 'test.txt' in testdata folder")
    if _, err := os.Stat("./testdata/test.txt"); os.IsNotExist(err) {
        t.Error("Expected file ./testdata/test.txt' to exist")
    }
}