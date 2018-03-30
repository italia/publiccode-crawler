package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/italia/developers-italia-backend/crawler"
	log "github.com/sirupsen/logrus"
)

var hostingFile = "hosting.yml"

func main() {
	var (
		version = flag.Bool("version", false, "prints version and exit")
	)

	flag.Parse()

	if *version {
		log.Info(Version)
		os.Exit(0)
	}

	data, err := ioutil.ReadFile(hostingFile)
	if err != nil {
		panic(fmt.Sprintf("error in reading %s file: %v", hostingFile, err))
	}

	hostings, err := crawler.ParseHostingFile(data)
	if err != nil {
		panic(fmt.Sprintf("error in parsing %s file: %v", hostingFile, err))
	}

	repositories := make(chan crawler.Repository, 100)
	for _, hosting := range hostings {
		go crawler.Process(hosting, repositories)
	}

	for {
		log.Info(<-repositories)
	}
}
