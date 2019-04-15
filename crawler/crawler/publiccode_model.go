package crawler

// softwareES represents a software record in Elasticsearch
type softwareES struct {
	FileRawURL            string            `json:"fileRawURL"`
	ID                    string            `json:"id"`
	CrawlTime             string            `json:"crawltime"`
	ItRiusoCodiceIPALabel string            `json:"it-riuso-codiceIPA-label"`
	Slug                  string            `json:"slug"`
	PublicCode            interface{}       `json:"publiccode"`
	VitalityScore         float64           `json:"vitalityScore"`
	VitalityDataChart     []int             `json:"vitalityDataChart"`
	OEmbedHTML            map[string]string `json:"oEmbedHTML"`
}
