package filestore

import (
	"fmt"
	"io"
	"os"
	"math/rand"
	"path/filepath"
	"time"
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
	rand.Seed(time.Now().UnixNano())
	fn := fmt.Sprintf("test%s.txt", randString(10))
	go func() {
		defer writer.Close()
		part, _ := writer.CreateFormFile("file", fn)
		content := "Test Add function"
		_, err := io.WriteString(part, content)
		if err != nil {
			t.Errorf("%v", err)
		}
	}()

	req := httptest.NewRequest("POST", "/", pr)
    req.Header.Add("Content-Type", writer.FormDataContentType())
	handler := func(w http.ResponseWriter, r *http.Request) {
		fs.Add(w, r)
	}
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	t.Logf("It should respond with an HTTP status code of 200")
    if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
        t.Errorf("Expected %d, received %d", 201, resp.StatusCode)
    }
    t.Logf("It should create a file named '%s' in testdata folder", fn)
    if _, err := os.Stat(filepath.Join("./testdata", fn)); os.IsNotExist(err) {
        t.Error("Expected file ./testdata/test.txt' to exist")
    }
}

func randString(n int) string {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}