package common

import (
	"strings"
)

// SplitFullName split a git FullName format to vendor and repo strings.
func SplitFullName(fullName string) (string, string) {
	s := strings.Split(fullName, "/")

	return s[0], s[1]
}
