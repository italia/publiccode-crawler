package jekyll

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"sort"

	"github.com/spf13/viper"

	"github.com/ghodss/yaml"
	"github.com/icza/dyno"
	"github.com/italia/developers-italia-backend/elastic"
	es "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

// software is used for parsing some fields of the software objects stored
// in Elasticsearch that are needed for computing additional information
// and for exporting variants and related software.
type software struct {
	ID         string `json:"id"`
	Slug       string `json:"slug"`
	PublicCode struct {
		URL         string   `json:"url"`
		Name        string   `json:"name"`
		IsBasedOn   []string `json:"isBasedOn"`
		Description map[string]struct {
			LocalisedName string   `json:"localisedName"`
			GenericName   string   `json:"genericName"`
			Features      []string `json:"features"`
			Screenshots   []string `json:"screenshots,omitempty"`
		} `json:"description"`
		It struct {
			Riuso struct {
				CodiceIPA string `json:"codiceIPA,omitempty"`
			} `json:"riuso"`
		} `json:"it"`
		Categories []string `json:"categories"`
		Legal      struct {
			RepoOwner string `json:"repoOwner,omitempty"`
		} `json:"legal,omitempty"`
	} `json:"publiccode"`

	// This is not populated from ES
	variants []software
}

// AllSoftwareYML generate the softwares.yml file
func AllSoftwareYML(filename string, numberOfSimilarSoftware, numberOfPopularCategories int, elasticClient *es.Client) error {
	log.Infof("Generating %s", filename)
	// Create file if not exists.
	if _, err := os.Stat(filename); os.IsExist(err) {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	// Open file.
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	// Extract all the softwares.
	query := elastic.NewBoolQuery("software")
	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		Pretty(true).                                       // pretty print request and response JSON
		From(0).Size(10000).                                // get first 10k elements. The limit can be changed in ES.
		Do(context.Background())                            // execute
	if err != nil {
		log.Error(err)
	}

	for _, hit := range searchResult.Hits.Hits {
		// hit.Source contains the raw JSON
		// We parse it into the first item of a slice, so that we can generate
		// YAML that looks like a single item and we can append it to the output
		// file as we go, without keeping all items in memory.
		full := make([]interface{}, 1)
		if err := json.Unmarshal(hit.Source, &full[0]); err != nil {
			log.Error(err)
		}

		// Let's parse the record again to get the fields we need for computing
		// additional information.
		var sw software
		if err := json.Unmarshal(hit.Source, &sw); err != nil {
			log.Error(err)
		}

		// Populate the output object with additional information
		dyno.Set(full[0], sw.findVariants(elasticClient), "oldVariant")
		dyno.Set(full[0], sw.variantsFeatures(), "oldFeatures")
		dyno.Set(full[0], sw.findRelated(numberOfSimilarSoftware, elasticClient), "relatedSoftwares")
		dyno.Set(full[0], sw.getPopularCategories(numberOfPopularCategories, elasticClient), "popularCategories")

		// Convert it to YAML
		yaml, err := yaml.Marshal(&full)
		if err != nil {
			log.Error(err)
		}

		// Append data to file.
		if _, err = f.WriteString(string(yaml)); err != nil {
			return err
		}
	}

	return err
}

// findVariants returns a list of variants of the given software.
func (sw *software) findVariants(elasticClient *es.Client) []software {
	query := elastic.NewBoolQuery("software")
	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		Pretty(true).                                       // pretty print request and response JSON
		From(0).Size(10000).                                // get first 10k elements. The limit can be changed in ES.
		Do(context.Background())                            // execute
	if err != nil {
		log.Error(err)
	}

	var sws []software
	for _, item := range searchResult.Each(reflect.TypeOf(*sw)) {
		i := item.(software)

		// TODO: this filtering logic should be moved to the ES query

		// skip identity
		if i.PublicCode.URL == sw.PublicCode.URL {
			continue
		}

		if funk.Contains(sw.PublicCode.IsBasedOn, i.PublicCode.URL) || funk.Contains(i.PublicCode.IsBasedOn, sw.PublicCode.URL) {
			sws = append(sws, i)
		}
	}
	return sws
}

// variantsFeatures returns features of variants that are not included in this one
func (sw *software) variantsFeatures() map[string][]string {
	diff := map[string][]string{} // "it" => [ feature, feature ... ]

	for _, lang := range []string{"en", "it"} {
		for _, variant := range sw.variants {
			for _, oldFeature := range variant.PublicCode.Description[lang].Features {
				if !funk.Contains(sw.PublicCode.Description[lang].Features, oldFeature) {
					diff[lang] = append(diff[lang], oldFeature)
				}
			}
		}
		diff[lang] = funk.UniqString(diff[lang])
	}

	return diff
}

// findRelated returns a list of similar software based on categories.
func (sw *software) findRelated(numberOfSimilarSoftware int, elasticClient *es.Client) []software {
	query := elastic.NewBoolQuery("software")
	for _, tag := range sw.PublicCode.Categories {
		query = query.Should(es.NewTermQuery("publiccode.categories", tag))
	}
	query = query.MustNot(es.NewTermsQuery("id", sw.ID))

	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		Sort("_score", false).                              // specify the sort order
		From(0).Size(numberOfSimilarSoftware).              // take documents from 0-numberOfSimilarSoftware
		Pretty(true).                                       // pretty print request and response JSON
		Do(context.Background())                            // execute
	if err != nil {
		log.Error(err)
	}

	var sws []software
	for _, item := range searchResult.Each(reflect.TypeOf(*sw)) {
		i := item.(software)
		sws = append(sws, i)
	}
	return sws
}

func (sw *software) getPopularCategories(number int, elasticClient *es.Client) []string {
	if len(sw.PublicCode.Categories) < number {
		return sw.PublicCode.Categories
	}

	// Extract all the documents. It should filter only the ones with isBaseOn=url.
	query := elastic.NewBoolQuery("software")
	searchResult, err := elasticClient.Search().
		Index(viper.GetString("ELASTIC_PUBLICCODE_INDEX")). // search in index "publiccode"
		Query(query).                                       // specify the query
		Pretty(true).                                       // pretty print request and response JSON
		From(0).Size(10000).                                // get first 10k elements. The limit can be changed in ES.
		Do(context.Background())                            // execute
	if err != nil {
		log.Error(err)
	}

	results := map[string]int{}

	// Range over the publiccodes in ES.
	for _, item := range searchResult.Each(reflect.TypeOf(*sw)) {
		i := item.(software)
		for _, v := range i.PublicCode.Categories {
			results[v]++
		}
	}

	// Order the map into a slice.
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range results {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// Populate the popularCategories slice with most popular categories.
	var popularCategories []string
	for n, kv := range ss {
		if n < number {
			break
		}
		popularCategories = append(popularCategories, kv.Key)
	}

	return popularCategories
}
