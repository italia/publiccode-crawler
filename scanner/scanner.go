package scanner

import (
	"errors"
	"net/url"

	"github.com/italia/publiccode-crawler/v4/common"
)

var ErrPubliccodeNotFound = errors.New("publiccode.yml not found")

// Scanner scans a single repository and emits a common.Repository on the
// repositories channel.
type Scanner interface {
	Scan(repoURL url.URL, publisher common.Publisher, repositories chan common.Repository) error
}
