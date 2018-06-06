package crawler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// ElasticClientFactory returns an elastic Client.
func ElasticClientFactory(URL, user, password string) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(URL),
		elastic.SetRetrier(NewESRetrier()),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(user, password),
		elastic.SetHealthcheckTimeoutStartup(60*time.Second))
	if err != nil {
		return nil, err
	}
	if elastic.IsConnErr(err) {
		log.Errorf("Elasticsearch connection problem: %v", err)
		return nil, err
	}

	return client, nil
}

// ElasticAliasUpdate update the Alias to the index.
func ElasticAliasUpdate(index, alias string, elasticClient *elastic.Client) error {
	// Remove old aliases.
	res, err := elasticClient.Aliases().Index("_all").Do(context.Background())
	if err != nil {
		return err
	}
	aliasService := elasticClient.Alias()
	indices := res.IndicesByAlias(alias)
	for _, name := range indices {
		log.Debugf("Remove alias from %s to %s", "publiccode", name)
		_, err := aliasService.Remove(name, alias).Do(context.Background())
		if err != nil {
			return err
		}

	}

	// Add an alias to the new index.
	log.Debugf("Add alias from %s to %s", index, alias)
	_, err = aliasService.Add(index, alias).Do(context.Background())

	return err
}

// ElasticRetrier implements the elastic interface that user can implement to intercept failed requests.
type ElasticRetrier struct {
	backoff elastic.Backoff
}

// NewESRetrier returns a new ElasticRetrier with Exponential Backoff waiting.
func NewESRetrier() *ElasticRetrier {
	return &ElasticRetrier{
		backoff: elastic.NewExponentialBackoff(10*time.Millisecond, 8*time.Second),
	}
}

// Retry is used in ElasticRetrier and returns the time to wait and if the retries should stop.
func (r *ElasticRetrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	log.Warn("Elasticsearch connection problem. Retry.")

	// Stop after 8 retries: ~2m.
	if retry >= 8 {
		return 0, false, errors.New("elasticsearch or network down")
	}

	// Let the backoff strategy decide how long to wait and whether to stop.
	wait, stop := r.backoff.Next(retry)
	return wait, stop, nil
}
