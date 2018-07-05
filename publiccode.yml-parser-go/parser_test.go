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
		// A complete and valid yml.
		{"tests/valid.yml", ""},
		// A complete and valid minimal yml.
		//{"tests/valid.minimal.yml", ""},

		// Missing mandatory fields.
		// {"tests/missing_publiccode-yaml-version.yml", "publiccode-yaml-version"},                 // Missing version.
		// {"tests/missing_name.yml", "name"},                                                       // Missing name.
		// {"tests/missing_legal_license.yml", "legal/license"},                                     // Missing legal/license.
		// {"tests/missing_legal_repoOwner.yml", "legal/repoOwner"},                                 // Missing legal/repoOwner.
		// {"tests/missing_localisation_availableLanguages.yml", "localisation/availableLanguages"}, // Missing localisation/availableLanguages.
		// {"tests/missing_localisation_localisationReady.yml", "localisation/localisationReady"},   // Missing localisation/localisationReady.
		// {"tests/missing_maintenance_contacts.yml", "maintenance/contacts"},                       // Missing maintenance/contacts.
		// {"tests/missing_maintenance_type.yml", "maintenance/type"},                               // Missing maintenance/type.
		// {"tests/missing_platforms.yml", "platforms"},                                             // Missing platforms.
		// {"tests/missing_releaseDate.yml", "releaseDate"},                                         // Missing releaseDate.
		// {"tests/missing_softwareType.yml", "softwareType"},                                       // Missing softwareType/type.
		// {"tests/missing_softwareVersion.yml", "softwareVersion"},                                 // Missing softwareVersion.
		// {"tests/missing_tags.yml", "tags"},                                                       // Missing tags.
		// {"tests/missing_url.yml", "url"},                                                         // Missing url.
	}

	for _, test := range testFiles {

		t.Run(test.errkey, func(t *testing.T) {
			// All tests are run in parallel with each other.
			t.Parallel()
			// Read data.
			data, err := ioutil.ReadFile(test.file)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Parse data into pc struct.
			var pc PublicCode
			err = Parse(data, &pc)

			//spew.Dump(pc.Description["eng"])

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
					t.Errorf("wrong error generated: %#v - instead of %s", e.Key, test.errkey)
				}
			}

		})
	}

}

// Test publiccode.yml remote files for key errors.
func TestDecodeValueErrorsRemote(t *testing.T) {
	BaseDir = "https://raw.githubusercontent.com/gith003/publiccode-org3/master/"

	testRemoteFiles := []struct {
		file   string
		errkey string
	}{
		// // A complete and valid REMOTE yml
		{"https://raw.githubusercontent.com/gith003/publiccode-org3/master/publiccode.yml", ""}, // Valid remote publiccode.yml.
		// //
		// // // A complete but invalid REMOTE yml
		// {"https://bitbucket.org/marco-capobussi/publiccode-example/raw/master/publiccode.yml-invalid", "publiccode-yaml-version"}, // Invalid remote publiccode.yml.
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
