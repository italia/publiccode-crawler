package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PaesslerAG/jsonpath"
)

// JSON fetches a JSON document and extracts repository URLs via a JSONPath expression.
type JSON struct {
	jsonPath string
}

func NewJSON(jsonPath string) *JSON {
	return &JSON{jsonPath: jsonPath}
}

func (c *JSON) Enumerate(ctx context.Context, u url.URL) ([]url.URL, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("json catalog: new request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("json catalog: GET %s: %w", u.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("json catalog: GET %s: status %d", u.String(), resp.StatusCode)
	}

	var data any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("json catalog: decode %s: %w", u.String(), err)
	}

	result, err := jsonpath.Get(c.jsonPath, data)
	if err != nil {
		return nil, fmt.Errorf("json catalog: jsonpath %q on %s: %w", c.jsonPath, u.String(), err)
	}

	var rawURLs []any
	switch v := result.(type) {
	case []any:
		rawURLs = v
	default:
		rawURLs = []any{v}
	}

	urls := make([]url.URL, 0, len(rawURLs))

	for _, raw := range rawURLs {
		s, ok := raw.(string)
		if !ok {
			continue
		}

		parsed, err := url.Parse(s)
		if err != nil {
			continue
		}

		urls = append(urls, *parsed)
	}

	return urls, nil
}
