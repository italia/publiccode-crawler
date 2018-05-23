package crawler

import (
	"context"
	"sync"

	"net/http"

	"github.com/italia/developers-italia-backend/httpclient"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/olivere/elastic"

	log "github.com/sirupsen/logrus"
)

// Repository is a single code repository.
type Repository struct {
	Name       string
	FileRawURL string
	Domain     Domain
	Headers    map[string]string
}

func ProcessRepositories(repositories chan Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	log.Debug("Repositories are going to be processed...")

	for repository := range repositories {
		wg.Add(1)
		go checkAvailability(repository, index, wg, elasticClient)
	}
}

func checkAvailability(repository Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	name := repository.Name
	fileRawUrl := repository.FileRawURL
	domain := repository.Domain
	headers := repository.Headers

	// Increment counter for the number of repositories processed.
	metrics.GetCounter("repository_processed", index).Inc()

	resp, err := httpclient.GetURL(fileRawUrl, headers)
	// If it's available and no error returned.
	if resp.Status.Code == http.StatusOK && err == nil {

		// Save to file.
		SaveToFile(domain, name, resp.Body, index)

		// Save to ES.
		SaveToES(domain, name, resp.Body, index, elasticClient)

		// Validate file.
		// TODO: uncomment these lines when mapping and File structure are ready for publiccode.
		// TODO: now validation is useless because we test on .gitignore file.
		// err := validateRemoteFile(resp.Body, fileRawUrl, index)
		// if err != nil {
		// 	log.Warn("Validator fails for: " + fileRawUrl)
		// 	log.Warn("Validator errors:" + err.Error())
		// }
	}

	// Defer waiting group close.
	wg.Done()
}

// Process delegates the work to single PA crawlers.
func ProcessPA(pa PA, domains []Domain, repositories chan Repository, index string, wg *sync.WaitGroup) {
	log.Debugf("ProcessPA: %s", pa.CodiceIPA)
	// range over repositories.
	for _, repository := range pa.Repositories {
		for _, domain := range domains {
			// if repository API is the domain
			if domain.ClientApi == repository.API {
				for _, org := range repository.Organizations {
					log.Debugf("ProcessPADomain: %s - API: %s", org, repository.API)
					wg.Add(1)
					ProcessPADomain(org, domain, repositories, index, wg)
				}
			}
		}
	}
	wg.Done()

}

func ProcessPADomain(org string, domain Domain, repositories chan Repository, index string, wg *sync.WaitGroup) {
	var url string

	// Starting URL.
	// TODO: refactoring
	if domain.ClientApi == "github" {
		url = domain.ApiBaseUrl + org + "/repos"
	}
	if domain.ClientApi == "gitlab" {
		url = domain.ApiBaseUrl + org
	}
	if domain.ClientApi == "bitbucket" {
		url = domain.ApiBaseUrl + org + "/?pagelen=100"
	}

	for {
		log.Debugf("processAndGetNextURL handler:%s", url)
		nextURL, err := domain.processAndGetNextURL(url, wg, repositories)
		if err != nil {
			log.Errorf("error reading %s repository list: %v. NextUrl: %v", url, err, nextURL)
			log.Errorf("Retry: %s", nextURL)
			nextURL = url
		}

		// If end is reached, nextUrl is empty.
		if nextURL == "" {
			log.Infof("Url: %s - is the last one.", url)

			// WaitingGroupd
			wg.Done()
			return
		}
		// Update url to nextURL.
		url = nextURL
	}
}

// // Process delegates the work to single domain crawlers.
// func ProcessDomain(domain Domain, redisClient *redis.Client, repositories chan Repository, index string, wg *sync.WaitGroup) {
// 	// Base starting URL.
// 	url := domain.URL
//
// 	for {
// 		// Set the value of nextURL on redis to domain.Index that describe the current execution.
// 		err := redisClient.HSet(domain.Id, url, index).Err()
// 		if err != nil {
// 			log.Error(err)
// 		}
//
// 		nextURL, err := domain.processAndGetNextURL(url, wg, repositories)
// 		if err != nil {
// 			log.Errorf("error reading %s repository list: %v. NextUrl: %v", url, err, nextURL)
// 			log.Errorf("Retry: %s", nextURL)
// 			nextURL = url
// 		}
// 		// If reached, the repository list was successfully retrieved.
// 		// Delete the repository url from redis.
// 		err = redisClient.HDel(domain.Id, url).Err()
// 		if err != nil {
// 			log.Error(err)
// 		}
//
// 		// If end is reached, nextUrl is empty.
// 		if nextURL == "" {
// 			log.Infof("Url: %s - is the last one.", url)
//
// 			// WaitingGroupd
// 			wg.Done()
// 			return
// 		}
// 		// Update url to nextURL.
// 		url = nextURL
// 	}
// }

// // validateRemoteFile validate the remote file
// func validateRemoteFile(data []byte, url, index string) error {
// 	fileName := viper.GetString("CRAWLED_FILENAME")
// 	// Parse data into pc struct and validate.
// 	baseURL := strings.TrimSuffix(url, fileName)
// 	// Set remore URL for remote validation (it will check files availability).
// 	publiccode.BaseDir = baseURL
// 	var pc publiccode.PublicCode
//
// 	err := publiccode.Parse(data, &pc)
//
// 	if err != nil {
// 		return err
// 	}
//
// 	metrics.GetCounter("repository_file_saved_valid", index).Inc()
// 	return err
//
// }
// WaitingLoop waits until all the goroutines counter is zero and close the repositories channel.
func WaitingLoop(repositories chan Repository, index string, wg *sync.WaitGroup, elasticClient *elastic.Client) {
	wg.Wait()

	// Remove old aliases.
	res, err := elasticClient.Aliases().Index("_all").Do(context.Background())
	if err != nil {
		panic(err)
	}
	aliasService := elasticClient.Alias()
	indices := res.IndicesByAlias("publiccode")
	for _, name := range indices {
		log.Debugf("Remove alias from %s to %s", "publiccode", name)
		aliasService.Remove(name, "publiccode").Do(context.Background())
	}

	// Add an alias to the new index.
	log.Debugf("Add alias from %s to %s", index, "publiccode")
	aliasService.Add(index, "publiccode").Do(context.Background())

	close(repositories)
}

//
// func addIndex(index string, elasticClient *elastic.Client) error {
// 	// Use the IndexExists service to check if a specified index exists.
// 	exists, err := elasticClient.IndexExists(index).Do(context.Background())
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		// Create a new index.
// 		// TODO: When mapping will be available: client.CreateIndex(index).BodyString(mapping).Do(ctx).
// 		_, err = elasticClient.CreateIndex(index).Do(context.Background())
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// // UpdateIndex return an old index if found on redis. A new index if no one is found.
// func UpdateIndex(domains []Domain, redisClient *redis.Client, elasticClient *elastic.Client) (string, error) {
// 	for i, _ := range domains {
// 		// Check if there is an URL that wasn't correctly retrieved.
// 		keys, err := redisClient.HKeys(domains[i].Id).Result()
// 		if err != nil {
// 			return "", err
// 		}
// 		// If some repository is left the index should remain the last one saved.
// 		for _, key := range keys {
// 			if redisClient.HGet(domains[i].Id, key).Val() != "" {
// 				index := redisClient.HGet(domains[i].Id, key).Val()
// 				addIndex(index, elasticClient)
// 				if err != nil {
// 					return "", err
// 				}
//
// 				return index, nil
// 			}
// 		}
// 	}

// 	// If reached there will be a new index.
// 	index := strconv.FormatInt(time.Now().Unix(), 10)
// 	err := addIndex(index, elasticClient)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return index, nil
// }
// // ProcessSingleRepository process a single repository given his url and domain.
// func ProcessSingleRepository(url string, domain Domain, repositories chan Repository) error {
//
// 	err := domain.processSingleRepo(url, repositories)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
