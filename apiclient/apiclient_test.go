package apiclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) APIClient {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	viper.Set("API_BASEURL", server.URL)
	viper.Set("API_BEARER_TOKEN", "token")

	return NewClient()
}

func patchBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()

	bytes, err := io.ReadAll(r.Body)
	assert.NoError(t, err)

	var fields map[string]any
	assert.NoError(t, json.Unmarshal(bytes, &fields))

	return fields
}

func TestGetSoftwareByURL_includesInactive(t *testing.T) {
	var query string

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery

		_, _ = w.Write([]byte(`{"data": [{"id": "id1", "active": false}], "links": {}}`))
	})

	software, err := client.GetSoftwareByURL("https://example.org/repo.git")

	assert.NoError(t, err)
	assert.Contains(t, query, "all=true")
	assert.Equal(t, "id1", software.ID)
	assert.False(t, software.Active)
}

func TestGetCatalogSoftwareByURL_includesInactive(t *testing.T) {
	var query string

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery

		_, _ = w.Write([]byte(`{"data": [{"id": "id1", "active": false}], "links": {}}`))
	})

	software, err := client.GetCatalogSoftwareByURL("catalog1", "https://example.org/repo.git")

	assert.NoError(t, err)
	assert.Contains(t, query, "all=true")
	assert.Equal(t, "id1", software.ID)
}

func TestPatchSoftware_sendsActive(t *testing.T) {
	var fields map[string]any

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		fields = patchBody(t, r)

		_, _ = w.Write([]byte(`{}`))
	})

	err := client.PatchSoftware("id1", "https://example.org/repo.git", nil, "yml", true)

	assert.NoError(t, err)
	assert.Equal(t, true, fields["active"])
}

func TestPatchCatalogSoftware_sendsActive(t *testing.T) {
	var fields map[string]any

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		fields = patchBody(t, r)

		_, _ = w.Write([]byte(`{}`))
	})

	err := client.PatchCatalogSoftware("catalog1", "id1", "https://example.org/repo.git", nil, "yml", true)

	assert.NoError(t, err)
	assert.Equal(t, true, fields["active"])
}
