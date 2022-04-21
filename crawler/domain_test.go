package crawler

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
)

// IsGithub returns "true" if the url can use Github API.
// TODO: complete
func TestParseDomainsFile(t *testing.T) {
	// Disablle log output for this function.
	log.SetOutput(ioutil.Discard)

	// Domains data.
	data := `- host: "gitlab.com"
	basic-auth:
		- ""
- host: "github.com"
	basic-auth:
		- ""
`

	// Domains struct.
	var domains []Domain
	domains = append(domains, Domain{
		Host: "gitlab.com",
	})
	domains = append(domains, Domain{
		Host: "github.com",
	})

	result, _ := parseDomainsFile([]byte(data))

	for i, domain := range domains {
		if domain.Host != domains[i].Host {
			t.Logf("Expected %s == %s.", result[i].Host, domains[i].Host)
			t.Fail()
		}

	}
}
