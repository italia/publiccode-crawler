package common

type Catalog struct {
	ID                  string
	Name                string
	PublishersNamespace string
	Sources             []CatalogSource
}
