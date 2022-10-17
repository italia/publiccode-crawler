package common

import (
	"crypto/sha1"
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Repository is a single code repository. FileRawURL contains the direct url to the raw file.
type Repository struct {
	Name         string
	URL          url.URL
	CanonicalURL url.URL
	FileRawURL   string
	GitBranch    string
	Publisher    Publisher
	Headers      map[string]string
}

// generateID generates a hash based on unique git repo URL.
func (repo *Repository) GenerateID() string {
	hash := sha1.New()
	_, err := hash.Write([]byte(repo.URL.String()))
	if err != nil {
		log.Errorf("Error generating the repository hash: %+v", err)
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// generateSlug generates a readable unique string based on repository name.
func (repo *Repository) GenerateSlug() string {
	vendorAndName := strings.Replace(repo.Name, "/", "-", -1)
	vendorAndName = strings.ReplaceAll(vendorAndName, ".", "_")

	if repo.Publisher.Id == "" {
		ID := repo.GenerateID()
		return fmt.Sprintf("%s-%s", vendorAndName, ID[0:6])
	}

	return fmt.Sprintf("%s-%s", repo.Publisher.Id, vendorAndName)
}
