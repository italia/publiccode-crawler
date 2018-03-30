package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"net/http"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var hostingFile = "hosting.yml"

func main() {
	var (
		version = flag.Bool("version", false, "prints version and exit")
	)

	flag.Parse()

	processedCounter := initPrometheus()

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

	processRepositories(repositories, processedCounter)
}

func initPrometheus() prometheus.Counter {
	processedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "repository_processed",
		Help: "Number of repository processed.",
	})
	err := prometheus.Register(processedCounter)
	if err != nil {
		log.Errorf("error in registering Prometheus handler: %v:", err)
	}

	go startMetricsServer()

	return processedCounter
}

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe("0.0.0.0:8081", nil)
	if err != nil {
		log.Warningf("monitoring endpoint non available: %v: ", err)
	}
}

func processRepositories(repositories chan crawler.Repository, processedCounter prometheus.Counter) {
	for repository := range repositories {
		log.Info(repository)
		processedCounter.Inc()
	}
}
