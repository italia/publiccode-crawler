package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(allCmd)
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Crawl publiccode.yml from hostings.",
	Long: `Start the crawler on every host written on hosting.yml file.
Beware! May take days to complete.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init Prometheus for metrics.
		processedCounter := initPrometheus()

		// Open and read hosting file list.
		hostingFile := "hosting.yml"
		data, err := ioutil.ReadFile(hostingFile)
		if err != nil {
			panic(fmt.Sprintf("error in reading %s file: %v", hostingFile, err))
		}
		// Parse hosting file list.
		hostings, err := crawler.ParseHostingFile(data)
		if err != nil {
			panic(fmt.Sprintf("error in parsing %s file: %v", hostingFile, err))
		}

		// Initiate a channel of repositories. TODO: check the limit. Possible Bottleneck.
		repositories := make(chan crawler.Repository)

		// For each host parsed from hosting, Process the repositories.
		for _, hosting := range hostings {
			go crawler.Process(hosting, repositories)
		}

		// Process the repositories in order to retrieve publiccode.yml.
		processRepositories(repositories, processedCounter)

	},
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
	// TODO: check the limit. Possible Bottleneck.
	ch := make(chan string, 100)
	counter := 0

	//Throttle requests.
	rate := time.Second / 1000 // 1 per second. TODO check rate limit (https://confluence.atlassian.com/bitbucket/rate-limits-668173227.html)
	throttle := time.Tick(rate)

	for repository := range repositories {
		// Throttle down the calls.
		<-throttle
		go checkAvailability(repository.Name, repository.URL, ch, processedCounter)

		// Comment: fmt.Println(counter)
		counter = counter + 1
	}
}

func checkAvailability(name, url string, ch chan<- string, processedCounter prometheus.Counter) {

	// Retrieve the url.
	response, err := http.Get(url)

	// If it's available and no error returned.
	if response.StatusCode == http.StatusOK && err == nil {
		log.Info("I FOUND ONE! IT'S: " + name + " at: " + url)

		// Retrieve the URL body.
		body, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()

		// Save the file.
		vendor, repo := splitFullName(name)
		fileName := "gitignore"
		go saveFile(vendor, repo, fileName, body)

		ch <- fmt.Sprintf("%s - FOUND IT! - %s", name, url)
	} else {
		ch <- fmt.Sprintf("%s - this one is bad :( - %s", name, url)
	}

	processedCounter.Inc()
}

func saveFile(vendor, repo, fileName string, data []byte) {

	path := filepath.Join("./data", vendor, repo)

	// MkdirAll will create all the folder path.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	err := ioutil.WriteFile(path+"/"+fileName, data, 0644)
	if err != nil {
		log.Info(err)
	}
}

func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}
