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
	metrics "github.com/italia/developers-italia-backend/metrics"
	publiccode "github.com/italia/developers-italia-backend/publiccode"
	"github.com/prometheus/client_golang/prometheus"
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
		processedCounter := metrics.PrometheusCounter("repository_processed", "Number of repository processed.")

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
		log.Debug("Loaded and parsed hosting.yml")

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

func processRepositories(repositories chan crawler.Repository, processedCounter prometheus.Counter) {
	log.Debug("Repositories are going to be processed...")
	// Throttle requests.
	// Time limits should be calibrated on more tests in order to avoid errors and bans.
	// 1/100 can perform a number of request < bitbucket limit.
	throttleRate := time.Second / 100
	throttle := time.Tick(throttleRate)

	for repository := range repositories {
		// Throttle down the calls.
		<-throttle
		go checkAvailability(repository.Name, repository.URL, repository.Headers, processedCounter)

	}

}

func checkAvailability(fullName, url string, headers map[string]string, processedCounter prometheus.Counter) {
	processedCounter.Inc()
	body, status, err := httpclient.GetURL(url, headers)
	// If it's available and no error returned.
	if status.StatusCode == http.StatusOK && err == nil {
		// Save the file.
		vendor, repo := splitFullName(fullName)
		fileName := os.Getenv("CRAWLED_FILENAME")
		saveFile(vendor, repo, fileName, body)

		// Parse data into pc struct and validate.
		baseURL := strings.TrimSuffix(url, "/publiccode.yml")
		publiccode.BaseDir = baseURL
		var pc publiccode.PublicCode

		err = publiccode.Parse(body, &pc)
		if err != nil {
			log.Warn("Invalid publiccode.yml: " + url)
			log.Warn("Errors:" + err.Error())
		}

	}
}

// saveFile save the choosen <file_name> in ./data/<vendor>/<repo>/<file_name>
func saveFile(vendor, repo, fileName string, data []byte) error {
	path := filepath.Join("./data", vendor, repo)

	// MkdirAll will create all the folder path, if not exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	err := ioutil.WriteFile(path+"/"+fileName, data, 0644)
	if err != nil {
		log.Error(err)
	}

	return err

}

// splitFullName split a git FullName format to vendor and repo strings.
func splitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")
	return s[0], s[1]
}
