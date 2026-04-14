package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PaesslerAG/jsonpath"
)

// JSONDriver fetches a JSON document representing a catalog and extracts
// repository URLs via a JSONPath expression.
type JSONDriver struct {
	jsonPath string
}

func NewJSONDriver(jsonPath string) *JSONDriver {
	return &JSONDriver{jsonPath: jsonPath}
}

func (c *JSONDriver) Enumerate(ctx context.Context, catalogURL url.URL) ([]url.URL, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, catalogURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("json catalog: new request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("json catalog: GET %s: %w", catalogURL.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("json catalog: GET %s: status %d", catalogURL.String(), resp.StatusCode)
	}

	var data any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("json catalog: decode %s: %w", catalogURL.String(), err)
	}

	result, err := jsonpath.Get(c.jsonPath, data)
	if err != nil {
		return nil, fmt.Errorf("json catalog: jsonpath %q on %s: %w", c.jsonPath, catalogURL.String(), err)
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
		rawURL, ok := raw.(string)
		if !ok {
			continue
		}

		parsed, err := url.Parse(rawURL)
		if err != nil {
			continue
		}

		urls = append(urls, *parsed)
	}

	return urls, nil
}
