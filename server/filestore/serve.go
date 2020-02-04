package filestore

import (
	"bufio"
	"fmt"
	"path/filepath"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"sort"
	"sync"

	"filestore/helper"
	"github.com/spf13/pflag"
	"github.com/sirupsen/logrus"
)

const (
	//DefaultPort for running filestore server
	DefaultPort int16 = 9090
)

// Config struct holds filestore server parameters
type Config struct {
	BindIP string
	BindHTTPPort int
	StoreDir string
	Logger  *logrus.Logger
}

// NewConfig return config initiated with default values
func NewConfig() *Config {
	c := Config{
		BindIP:          "0.0.0.0",
		BindHTTPPort:    int(DefaultPort),
		StoreDir: "",
		Logger: helper.NewLogger("filestore"),
	}
	return &c
}

//BindFlags attaches the  pflag flagset to the current config
func (c *Config) BindFlags(fs *pflag.FlagSet) {
	home, _ := os.UserHomeDir()
	fs.StringVar(&c.BindIP, "bind-ip", c.BindIP, "IP fileStore server will listen to")
	fs.IntVar(&c.BindHTTPPort, "bind-http-port", c.BindHTTPPort, "HTTP Port fileStore server will listen to")
	fs.StringVar(&c.StoreDir, "store-dir", filepath.Join(home,"store"), "filestore storage dir")
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of filestore:\n")
		fs.PrintDefaults()
	}
}
// FileStore is the root data object of the filestore
type FileStore struct {
	Logger *logrus.Logger
	StoreDir string
	BindHTTPAddress string
}

// init creates the store if it doesnt exist
func (fs *FileStore) init() {
	if _, err := os.Stat(fs.StoreDir); os.IsNotExist(err) {
			err = os.MkdirAll(fs.StoreDir, 0755)
			if err != nil {
				fs.Logger.Fatalf("Could not create file store: %s", err)
			}
	}
}

// NewFileStore creates a new Client
func NewFileStore(c *Config) *FileStore {
	HTTPaddress := net.JoinHostPort(c.BindIP, strconv.Itoa(c.BindHTTPPort))
	var fs = FileStore{
		BindHTTPAddress: HTTPaddress,
		Logger: c.Logger,
		StoreDir: c.StoreDir,
	}
	fs.init()
	return &fs
}

// Run start a filestore server and make it listen for incoming connections
func (fs *FileStore) Run() {
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		fs.Add(w, r)
	})
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		fs.List(w, r)
	})
	http.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		fs.Remove(w, r)
	})
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		fs.Update(w, r)
	})
	http.HandleFunc("/freqwords", func(w http.ResponseWriter, r *http.Request) {
		fs.FreqWords(w, r)
	})
	http.HandleFunc("/countwords", func(w http.ResponseWriter, r *http.Request) {
		fs.CountWords(w, r)
	})
	HTTPServer := &http.Server{Addr: fs.BindHTTPAddress, Handler: nil}
	if err := HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fs.Logger.Fatalf("Could not listen on %s: %v\n", fs.BindHTTPAddress, err)
	}
}

// Add adds files to the store
func (fs *FileStore) Add(w http.ResponseWriter, r *http.Request) {
	fs.Logger.Infof("Adding multipart files to the store")
	reader, err := r.MultipartReader()
	if err != nil {
		fs.Logger.Fatalf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FileName() == "" {
			continue
		}
		fs.Logger.Infof("checking if file exist in the store")
		if _, err = os.Stat(filepath.Join(fs.StoreDir, part.FileName())); !os.IsNotExist(err) {
			http.Error(w, "File already exist", http.StatusConflict)
			return
		}
		fs.Logger.Infof("Adding file %s to the store", part.FileName())
		dst, err := os.Create(filepath.Join(fs.StoreDir, part.FileName()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// List lists files in the store
func (fs *FileStore) List(w http.ResponseWriter, r *http.Request) {
	fs.Logger.Infof("Listing files in the store")
	files, err := ioutil.ReadDir(fs.StoreDir)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }
    for _, file := range files {
		_, err = io.WriteString(w, fmt.Sprintf("%s\n", file.Name()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Remove removes files from the store
func (fs *FileStore) Remove(w http.ResponseWriter, r *http.Request) {
	fs.Logger.Infof("Removing file from the store")
	fileName := r.FormValue("file")
	fs.Logger.Infof("Removing file name %s", fileName)
	err := os.Remove(filepath.Join(fs.StoreDir, fileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Update updates a file in the store
func (fs *FileStore) Update(w http.ResponseWriter, r *http.Request) {
	fs.Logger.Infof("Updating file in the store")
	reader, err := r.MultipartReader()
	if err != nil {
		fs.Logger.Fatalf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	part, err := reader.NextPart()
	if err != nil {
		fs.Logger.Fatalf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fs.Logger.Infof("Updating file %s",part.FileName())
	dst, err := os.OpenFile(filepath.Join(fs.StoreDir,part.FileName()), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, part); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// FreqWords return most frequent words
func (fs *FileStore) FreqWords(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	order := queryValues.Get("order")
	limit, err := strconv.Atoi(queryValues.Get("limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fs.Logger.Infof("Computing most %s frequent words in %s ordering", queryValues.Get("limit"), queryValues.Get("order"))
	result, err := searchInDir(fs.StoreDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type keyval struct {
		key string
		val int
	}
	var sortedRes []keyval
    for k, v := range result {
        sortedRes = append(sortedRes, keyval{k, v})
	}

	sort.Slice(sortedRes, func(i, j int) bool {
		if order == "dsc" {
			return sortedRes[i].val > sortedRes[j].val
		} 
		return sortedRes[i].val < sortedRes[j].val
	})
	for rank := 0; (rank < limit && rank < len(sortedRes)); rank++ {
        word := sortedRes[rank].key
		freq := sortedRes[rank].val
		_, err = io.WriteString(w,fmt.Sprintf("%3d %s\n", freq, word))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
    }
}

// CountWords counts words in the store
func (fs *FileStore) CountWords(w http.ResponseWriter, r *http.Request) {
	fs.Logger.Infof("Counting words in the store")
	result, err := countInDir(fs.StoreDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.WriteString(w, fmt.Sprintf("%3d\n", result))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


// searchInDir return a map of words and their occurence inside a folder
func searchInDir(dir string) (map[string]int, error) {
	SuperResult := make(map[string]int)
	filelist, err := ioutil.ReadDir(dir)
	if err != nil {
		helper.NewLogger("filestore").Fatalf("%v", err)
		return nil, err
	}
	wg := &sync.WaitGroup{}
	resultChan := make(chan map[string]int)
	for _, fileinfo := range filelist {
		wg.Add(1)
		go searchInFile(filepath.Join(dir,fileinfo.Name()), resultChan, wg)
	}
	go func() {   
		wg.Wait()
		close(resultChan)
	}()
	for m := range resultChan {
		for k, v := range m {
			SuperResult[k] = SuperResult[k] + v
		}
	}
	return SuperResult, nil
}

// searchInFile scan a file and build a map of string occurence
func searchInFile(path string, resultChan chan map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(path)
	if err != nil {
		helper.NewLogger("filestore").Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	result := make(map[string]int)
	for scanner.Scan() {
		result[scanner.Text()]++
	}
	resultChan <- result
}

// countInDir counts the words in a folder
func countInDir(dir string) (int, error) {
	filelist, err := ioutil.ReadDir(dir)
	if err != nil {
		helper.NewLogger("filestore").Fatalf("%v", err)
		return 0, err
	}
	wg := &sync.WaitGroup{}
	resultChan := make(chan int)
	for _, fileinfo := range filelist {
		wg.Add(1)
		go countInFile(filepath.Join(dir,fileinfo.Name()), resultChan, wg)
	}
	go func() {   
		wg.Wait()
		close(resultChan)
	}()

	SuperCount := 0
	for i := range resultChan {
		SuperCount = SuperCount + i
	}
	return SuperCount, nil
}

// countInFile counts the words in a file
func countInFile(path string, resultChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(path)
	if err != nil {
		helper.NewLogger("filestore").Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	result := 0
	for scanner.Scan() {
		result++
	}
	resultChan <- result
}