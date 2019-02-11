package crawler

// SoftwareES represents a software record in Elasticsearch
type SoftwareES struct {
	FileRawURL            string            `json:"fileRawURL"`
	ID                    string            `json:"id"`
	CrawlTime             string            `json:"crawltime"`
	ItRiusoCodiceIPALabel string            `json:"it-riuso-codiceIPA-label"`
	PublicCode            interface{}       `json:"publiccode"`
	VitalityScore         float64           `json:"vitalityScore"`
	VitalityDataChart     []int             `json:"vitalityDataChart"`
	OEmbedHTML            map[string]string `json:"oEmbedHTML"`
}
