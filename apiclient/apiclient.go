package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/italia/publiccode-crawler/v3/common"
	internalUrl "github.com/italia/publiccode-crawler/v3/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ApiClient struct {
	baseURL         string
	retryableClient *http.Client
	token           string
}

type Links struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

type PublishersPaginated struct {
	Data  []Publisher `json:"data"`
	Links Links       `json:"links"`
}

type SoftwarePaginated struct {
	Data  []Software `json:"data"`
	Links Links      `json:"links"`
}

type Publisher struct {
	ID            string        `json:"id"`
	AlternativeID string        `json:"alternativeId"`
	Email         string        `json:"email"`
	Description   string        `json:"description"`
	CodeHostings  []CodeHosting `json:"codeHosting"`
	Active        bool          `json:"active"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

type CodeHosting struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Group     bool      `json:"group"`
	URL       string    `json:"url"`
}

type Software struct {
	ID            string    `json:"id"`
	URL           string    `json:"url"`
	Aliases       []string  `json:"aliases"`
	PubliccodeYml string    `json:"publiccodeYml"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewClient() ApiClient {
	retryClient := retryablehttp.NewClient()

	retryClient.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	retryClient.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		logrus.WithFields(logrus.Fields{
			"proto":    req.Proto,
			"endpoint": req.URL.String(),
			"attempt":  attempt,
		}).Trace("Sending request to API")
	}

	// Create a standard *http.Client with the retryablehttp.Client as the transport

	return ApiClient{
		baseURL:         viper.GetString("API_BASEURL"),
		retryableClient: retryClient.StandardClient(),
		token:           "Bearer " + viper.GetString("API_BEARER_TOKEN"),
	}
}

func (clt ApiClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return clt.retryableClient.Do(req)
}

func (clt ApiClient) Post(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", clt.token)
	req.Header.Add("Content-Type", "application/json")

	return clt.retryableClient.Do(req)
}

func (clt ApiClient) Patch(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPatch,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", clt.token)
	req.Header.Add("Content-Type", "application/merge-patch+json")

	return clt.retryableClient.Do(req)
}

// GetPublishers returns a slice with all the publishers from the API and
// any error encountered.
func (c ApiClient) GetPublishers() ([]common.Publisher, error) {
	var publishersResponse *PublishersPaginated
	var publishers []common.Publisher

	pageAfter := ""

page:
	reqUrl := joinPath(c.baseURL, "/publishers") + pageAfter

	res, err := c.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("can't get publishers %s: %w", reqUrl, err)
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("can't get publishers %s: HTTP status %s", reqUrl, res.Status)
	}

	publishersResponse = &PublishersPaginated{}
	err = json.NewDecoder(res.Body).Decode(&publishersResponse)
	if err != nil {
		return nil, fmt.Errorf("can't parse GET %s response: %w", reqUrl, err)
	}

	for _, p := range publishersResponse.Data {
		var groups, repos []internalUrl.URL

		for _, codeHosting := range p.CodeHostings {
			u, err := url.Parse(codeHosting.URL)
			if err != nil {
				return nil, fmt.Errorf("can't parse GET %s response: %w", reqUrl, err)
			}

			if codeHosting.Group {
				groups = append(groups, (internalUrl.URL)(*u))
			} else {
				repos = append(repos, (internalUrl.URL)(*u))
			}
		}

		id := p.ID

		// Let's give precedence to the alternativeId. It's usually set by
		// at Publisher creation and it's supposed to be more representative of the
		// Publisher than the automatically generated UUID, since it's explicitly
		// set with the API by the creator.
		//
		// This way, we also take a minimalist approach to Publisher concept in the crawler,
		// having just one id.
		if p.AlternativeID != "" {
			id = p.ID
		}
		publishers = append(publishers, common.Publisher{
			Id:            id,
			Name:          fmt.Sprintf("%s %s", p.Description, p.Email),
			Organizations: groups,
			Repositories:  repos,
		})
	}

	if publishersResponse.Links.Next != "" {
		pageAfter = publishersResponse.Links.Next

		goto page
	}

	return publishers, nil
}

// GetSoftwareByURL returns the software matching the given repo URL and
// any error encountered.
func (c ApiClient) GetSoftwareByURL(url string) (*Software, error) {
	var softwareResponse SoftwarePaginated

	res, err := c.retryableClient.Get(joinPath(c.baseURL, "/software") + "?url=" + url)
	if err != nil {
		return nil, fmt.Errorf("can't GET /software?url=%s: %w", url, err)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&softwareResponse)
	if err != nil {
		return nil, fmt.Errorf("can't parse GET /software?url=%s response: %w", url, err)
	}

	if len(softwareResponse.Data) > 0 {
		return &softwareResponse.Data[0], nil
	}

	return nil, nil
}

// PostSoftware creates a new software resource with the given fields and returns
// a Software struct or any error encountered.
func (c ApiClient) PostSoftware(url string, aliases []string, publiccodeYml string, active bool) (*Software, error) {
	body, err := json.Marshal(map[string]interface{}{
		"publiccodeYml": publiccodeYml,
		"url":           url,
		"aliases":       aliases,
		"active":        active,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create software: %w", err)
	}

	res, err := c.Post(joinPath(c.baseURL, "/software"), body)
	if err != nil {
		return nil, fmt.Errorf("can't create software: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("can't create software: API replied with HTTP %s", res.Status)
	}

	postSoftwareResponse := &Software{}
	err = json.NewDecoder(res.Body).Decode(&postSoftwareResponse)
	if err != nil {
		return nil, fmt.Errorf("can't parse POST /software (for %s) response: %w", url, err)
	}

	return postSoftwareResponse, nil
}

// PatchSoftware updates a software resource with the given fields and returns
// an http.Response and any error encountered.
func (c ApiClient) PatchSoftware(id string, url string, aliases []string, publiccodeYml string) (*http.Response, error) {
	body, err := json.Marshal(map[string]interface{}{
		"publiccodeYml": publiccodeYml,
		"url":           url,
		"aliases":       aliases,
	})
	if err != nil {
		return nil, fmt.Errorf("can't update software: %w", err)
	}

	res, err := c.Patch(joinPath(c.baseURL, "/software/"+id), body)
	if err != nil {
		return res, fmt.Errorf("can't update software: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return res, fmt.Errorf("can't update software: API replied with HTTP %s", res.Status)
	}

	return res, nil
}

// PostSoftwareLog creates a new software log with the given fields and returns
// an http.Response and any error encountered.
func (c ApiClient) PostSoftwareLog(softwareId string, message string) (*http.Response, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"message": message,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create log: %w", err)
	}

	res, err := c.Post(joinPath(c.baseURL, "/software/", softwareId, "logs"), payload)
	if err != nil {
		return res, fmt.Errorf("can't create software log: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return res, fmt.Errorf("can't create software log: API replied with HTTP %s", res.Status)
	}

	return res, nil
}

// PostLog creates a new log with the given message and returns an http.Response
// and any error encountered.
func (c ApiClient) PostLog(message string) (*http.Response, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"message": message,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create log: %w", err)
	}

	res, err := c.Post(joinPath(c.baseURL, "/logs"), payload)
	if err != nil {
		return res, fmt.Errorf("can't create log: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return res, fmt.Errorf("can't create log: API replied with HTTP %s", res.Status)
	}

	return res, nil
}

func joinPath(base string, paths ...string) string {
	u, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range paths {
		u.Path = path.Join(u.Path, p)
	}

	return u.String()
}
