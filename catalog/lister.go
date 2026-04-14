package catalog

import (
	"net/url"

	"github.com/italia/publiccode-crawler/v4/common"
)

// Lister enumerates all repositories under a group/organization URL and emits
// a common.Repository for each one on the repositories channel.
type Lister interface {
	List(groupURL url.URL, publisher common.Publisher, repositories chan common.Repository) error
}
