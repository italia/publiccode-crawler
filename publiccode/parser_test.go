package publiccode

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// Test publiccode.yml local files for key errors.
func TestDecodeValueErrors(t *testing.T) {
	BaseDir = ""

	testFiles := []struct {
		file   string
		errkey string
	}{
		// A complete and valid yml
		{"tests/valid.yml", ""}, // Valid yml.

		// // Version
		// {"tests/invalid_version.yml", "version"}, // Invalid version.
		//
		// // Url
		// {"tests/invalid_url_schema.yml", "url"},      // Missing schema.
		// {"tests/invalid_url_404notfound.yml", "url"}, // 404 not found.
		//
		// // UpstreamURL
		// {"tests/valid_upstream-url_missing.yml", ""},                   // Valid. Missing non-mandatory.
		// {"tests/invalid_upstream-url_schema.yml", "upstream-url"},      // Missing schema.
		// {"tests/invalid_upstream-url_404notfound.yml", "upstream-url"}, // 404 not found.
		//
		// //Legal
		// {"tests/valid_legal_missing.yml", ""},                              // Valid. Missing non-mandatory.
		// {"tests/invalid_legal-repo-owner_missing.yml", "legal/repo-owner"}, // Missing legal/repo-owner.
		// {"tests/invalid_legal-license_missing.yml", "legal/license"},       // Missing legal/license.
		// {"tests/invalid_legal-license_nospdxlicense.yml", "legal/license"}, // Non-SPDX license.

	}

	for _, test := range testFiles {
		t.Run(test.errkey, func(t *testing.T) {

			// Read data.
			data, err := ioutil.ReadFile(test.file)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Parse data into pc struct.
			var pc PublicCode
			err = Parse(data, &pc)

			if test.errkey == "" && err != nil {
				t.Error("unexpected error:\n", err)
			} else if test.errkey != "" && err == nil {
				t.Error("error not generated:\n", test.file)
			} else if test.errkey != "" && err != nil {
				if multi, ok := err.(ErrorParseMulti); !ok {
					panic(err)
				} else if len(multi) != 1 {
					t.Errorf("too many errors generated: %#v", multi)
				} else if e, ok := multi[0].(ErrorInvalidValue); !ok || e.Key != test.errkey {
					t.Errorf("wrong error generated: %#v - instead of %s", e, test.errkey)
				}
			}
		})
	}
}

// Test publiccode.yml remote files for key errors.
func TestDecodeValueErrorsRemote(t *testing.T) {
	BaseDir = "https://bitbucket.org/marco-capobussi/publiccode-example/raw/master/"

	testRemoteFiles := []struct {
		file   string
		errkey string
	}{
		// A complete and valid REMOTE yml
		{"https://bitbucket.org/marco-capobussi/publiccode-example/raw/master/publiccode.yml", ""}, // Valid remote publiccode.yml.

		// A complete but invalid REMOTE yml
		{"https://bitbucket.org/marco-capobussi/publiccode-example/raw/master/publiccode.yml-invalid", "description/logo"}, // Invalid remote publiccode.yml.
	}

	for _, test := range testRemoteFiles {
		t.Run(test.errkey, func(t *testing.T) {

			// Read data.
			resp, err := http.Get(test.file)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Parse data into pc struct.
			var pc PublicCode
			err = Parse(data, &pc)

			if test.errkey == "" && err != nil {
				t.Error("unexpected error:\n", err)
			} else if test.errkey != "" && err == nil {
				t.Error("error not generated:\n", test.file)
			} else if test.errkey != "" && err != nil {
				if multi, ok := err.(ErrorParseMulti); !ok {
					panic(err)
				} else if len(multi) != 1 {
					t.Errorf("too many errors generated: %#v", multi)
				} else if e, ok := multi[0].(ErrorInvalidValue); !ok || e.Key != test.errkey {
					t.Errorf("wrong error generated: %#v - instead of %s", e, test.errkey)
				}
			}
		})
	}
}
