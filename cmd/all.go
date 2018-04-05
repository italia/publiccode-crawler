package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostingFile = "hosting.yml"

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

		// Open hosting file list.
		data, err := ioutil.ReadFile(hostingFile)
		if err != nil {
			panic(fmt.Sprintf("error in reading %s file: %v", hostingFile, err))
		}
		// Parse hosting file list
		hostings, err := crawler.ParseHostingFile(data)
		if err != nil {
			panic(fmt.Sprintf("error in parsing %s file: %v", hostingFile, err))
		}

		// For each host parsed from hosting, Process the repositories
		repositories := make(chan crawler.Repository, 100)
		for _, hosting := range hostings {
			go crawler.Process(hosting, repositories)
		}

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

	ch := make(chan string, 100)
	counter := 0

	//throttle requests
	rate := time.Second // 1 per second. TODO check rate limit (https://confluence.atlassian.com/bitbucket/rate-limits-668173227.html)
	throttle := time.Tick(rate)

	for repository := range repositories {
		//log.Info(repository)

		<-throttle
		go checkAvailability(repository.Name, repository.URL, ch, processedCounter)

		// slow version
		// response, err := http.Get(repository.URL)
		// if response.StatusCode == http.StatusOK && err == nil {
		// 	log.Info("I FOUND ONE! IT'S: " + repository.Name + " at: " + repository.URL)
		// } else {
		// 	log.Info("this one is bad :( " + repository.URL)
		// }
		fmt.Println(counter)
		counter = counter + 1
	}
}

func checkAvailability(name, url string, ch chan<- string, processedCounter prometheus.Counter) {

	response, err := http.Get(url)
	if response.StatusCode == http.StatusOK && err == nil {
		log.Info("I FOUND ONE! IT'S: " + name + " at: " + url)

		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		err := ioutil.WriteFile("data/"+name, body, 0644)
		if err != nil {
			log.Info("Error")
		}

		ch <- fmt.Sprintf("%s - FOUND IT! - %s", name, url)
	} else {
		ch <- fmt.Sprintf("%s - this one is bad :( - %s", name, url)
	}

	processedCounter.Inc()
}
