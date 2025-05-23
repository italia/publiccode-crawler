package internal

import (
	"net/url"
)

type URL url.URL

// UnmarshalYAML implements the yaml.Unmarshaler interface for URLs.
func (u *URL) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	urlp, err := url.Parse(s)
	if err != nil {
		return err
	}

	*u = (URL)(*urlp)

	return nil
}

func (u URL) MarshalYAML() (interface{}, error) {
	return u.String(), nil
}

func (u *URL) String() string {
	return (*url.URL)(u).String()
}
