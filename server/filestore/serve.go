package filestore

import (
	"fmt"
	"path/filepath"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"

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
		fs.Logger.Infof("Adding file %s to the store", part.FileName())
		dst, err := os.Create(filepath.Join(fs.StoreDir, part.FileName()))
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
		w.Write([]byte(fmt.Sprintf("%s\n", file.Name())))
	
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
