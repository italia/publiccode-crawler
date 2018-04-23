package publiccode

import (
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

// checkEmail tells whether email is well formatted.
// In general an email is valid if compile the regex: ^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$
func (p *parser) checkEmail(key string, fn string) error {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(fn) {
		return newErrorInvalidValue(key, "invalid email: %v", fn)
	}

	return nil
}

// checkUrl tells whether the URL resource is well formatted and reachable and return it as *url.URL.
// An URL resource is well formatted if it's' a valid URL and the scheme is not empty.
// An URL resource is reachable if returns an http Status = 200 OK.
func (p *parser) checkUrl(key string, value string) (*url.URL, error) {
	u, err := url.Parse(value)
	if err != nil {
		return nil, newErrorInvalidValue(key, "not a valid URL: %s", value)
	}
	if u.Scheme == "" {
		return nil, newErrorInvalidValue(key, "missing URL scheme: %s", value)
	}
	r, err := http.Get(value)
	if err != nil {
		return nil, newErrorInvalidValue(key, "Http.get failed for: %s", value)
	}
	if r.StatusCode != 200 {
		return nil, newErrorInvalidValue(key, "URL is unreachable: %s", value)
	}

	return u, nil
}

// checkFile tells whether the file resource exists and return it.
func (p *parser) checkFile(key string, fn string) (string, error) {
	if BaseDir == "" {
		if _, err := os.Stat(fn); err != nil {
			return "", newErrorInvalidValue(key, "file does not exist: %v", fn)
		}
	} else {
		//Remote bitbucket
		_, err := p.checkUrl(key, BaseDir+fn)

		//_, err := p.checkUrl(key, "https://bitbucket.org/marco-capobussi/publiccode-example/raw/master/"+fn)
		if err != nil {
			return "", newErrorInvalidValue(key, "file does not exist on remote: %v", BaseDir+fn)
		}
	}
	return fn, nil
}

// checkDate tells whether the string in input is a date in the
// format "YYYY-MM-DD", which is one of the ISO8601 allowed encoding, and return it as time.Time.
func (p *parser) checkDate(key string, value string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", value); err != nil {
		return t, newErrorInvalidValue(key, "cannot parse date: %v", err)
	} else {
		return t, nil
	}
}

// checkImage tells whether the string in a valid image. It also checks if the file exists.
func (p *parser) checkImage(key string, value string) error {
	// Based on https://github.com/italia/publiccode.yml/blob/master/schema.md#key-descriptionlogo
	//TODO: check two version of the Logo
	//TODO: check extensions and image size
	//TODO: raster should exists iff vector does not exists

	if _, err := p.checkFile(key, value); err != nil {
		return err
	}

	return nil
}
