package crawler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/olivere/elastic"
	"github.com/prometheus/common/log"
)

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
		log.Error("Elasticsearch connection problem: %v", err)
		return nil, err
	}

	return client, nil
}

type ElasticRetrier struct {
	backoff elastic.Backoff
}

func NewESRetrier() *ElasticRetrier {
	return &ElasticRetrier{
		backoff: elastic.NewExponentialBackoff(10*time.Millisecond, 8*time.Second),
	}
}

func (r *ElasticRetrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	log.Warn("Elasticsearch connection problem. Retry.")

	// Stop after 8 retries: 2m
	if retry >= 8 {
		return 0, false, errors.New("elasticsearch or network down")
	}

	// Let the backoff strategy decide how long to wait and whether to stop
	wait, stop := r.backoff.Next(retry)
	return wait, stop, nil
}
