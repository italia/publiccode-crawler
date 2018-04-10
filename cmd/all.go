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
	"github.com/italia/developers-italia-backend/httpclient"
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

		// Initiate a channel of repositories.
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
	channelCapacity := 100
	ch := make(chan string, channelCapacity)
	// Throttle requests.
	// Time limits should be calibrated on more tests in order to avoid errors and bans.
	// 1/100 can perform a number of request < bitbucket limit.
	rate := time.Second / 100
	throttle := time.Tick(rate)

	for repository := range repositories {
		// Throttle down the calls.
		<-throttle
		go checkAvailability(repository.Name, repository.URL, ch, processedCounter)

	}
}

func checkAvailability(name, url string, ch chan<- string, processedCounter prometheus.Counter) {
	// Retrieve the url.
	body, status, err := httpclient.GetURL(url)

	// If it's available and no error returned.
	if status.StatusCode == http.StatusOK && err == nil {
		// Save the file.
		vendor, repo := splitFullName(name)
		fileName := "gitignore"
		go saveFile(vendor, repo, fileName, body)

		ch <- fmt.Sprintf("%s - hit - %s", name, url)
	} else {
		ch <- fmt.Sprintf("%s - miss - %s", name, url)
	}

	processedCounter.Inc()
}

// saveFile save the choosen <file_name> in ./data/<vendor>/<repo>/<file_name>
func saveFile(vendor, repo, fileName string, data []byte) {
	path := filepath.Join("./data", vendor, repo)

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	err := ioutil.WriteFile(path+"/"+fileName, data, 0644)
	if err != nil {
		log.Info(err)
	}
}

// splitFullName split a git FullName format to vendor and repo strings.
func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}
