package main

import (
	"filestore/helper"
	"filestore/server/filestore"
	"github.com/spf13/pflag"
)

func main() {
	// Binding flags for the command line and allow parsing
	logger := helper.NewLogger("filestore")
	config := filestore.NewConfig()
	config.BindFlags(pflag.CommandLine)
	pflag.Parse()


	fs := filestore.NewFileStore(config)
	logger.Infof("Starting the filestore server")
	fs.Run()
}