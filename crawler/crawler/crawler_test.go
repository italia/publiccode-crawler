package crawler

import (
	"io/ioutil"
	"testing"

	publiccode "github.com/italia/publiccode-parser-go"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

var whitelist string = `
-
  name: testit
  repos: 
    - "https://github.com/test/testrepo"
-
  codice-iPA:
  name: testit
  repos: 
    - "https://github.com/test/testrepo"
-
  codice-iPA: test
  name: testit
  repos: 
    - "https://github.com/test/testrepo"
`

// validateRemoteFile will parse and validate
// crawled publiccode.yml. It will thrown errors
// if parse fails and if IPA code mismatch between
// whithelist file and publiccode itself
func TestIPAMatch(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	var pas []PA
	var parser publiccode.Parser
	err := yaml.Unmarshal([]byte(whitelist), &pas)
	if err != nil {
		t.Errorf("error on unmarsalling whitelist %s", err)
	}

	// should not throw error codiceIPA key is equal
	// on both sides
	for _, pa := range pas {
		parser.PublicCode.It.Riuso.CodiceIPA = pa.CodiceIPA
		err = validateFile(pa, parser, "")
		if err != nil {
			t.Errorf("error comparing IPA codes %s", err)
		}
	}

	// it should thowns errors since they always mismatch
	for _, pa := range pas {
		parser.PublicCode.It.Riuso.CodiceIPA = pa.CodiceIPA + "x"
		err = validateFile(pa, parser, "")
		if err == nil {
			t.Errorf("error comparing IPA codes %v", err)
		}
	}
}
