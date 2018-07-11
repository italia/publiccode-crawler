package crawler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
)

// ElasticClientFactory returns an elastic Client.
func ElasticClientFactory(URL, user, password string) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(URL),
		elastic.SetRetrier(NewESRetrier()),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(user, password),
		elastic.SetHealthcheckTimeoutStartup(60*time.Second))
	if err != nil {
		return nil, err
	}
	if elastic.IsConnErr(err) {
		log.Errorf("Elasticsearch connection problem: %v", err)
		return nil, err
	}

	return client, nil
}

// ElasticIndexMapping adds (if not exists) the mapping for the crawler data in ES.
func ElasticIndexMapping(index string, elasticClient *elastic.Client) error {
	const (
		// Elasticsearch mapping for publiccode. Check elasticsearch/mappings/.
		mapping = `{
		"mappings": {
			"software": {
				"dynamic_templates": [
					{
						"description": {
							"path_match": "description.*",
							"mapping": {
								"type": "object",
								"properties": {
									"localisedName": {
										"type": "text"
									},
									"genericName": {
										"type": "text",
										"fields": {
											"keyword": { "type": "keyword", "ignore_above": 256 }
										}
									},
									"shortDescription": {
										"type": "text"
									},
									"longDescription": {
										"type": "text"
									},
									"documentation": {
										"type": "text",
										"index": false
									},
									"apiDocumentation": {
										"type": "text",
										"index": false
									},
									"featureList": {
										"type": "keyword"
									},
									"freeTags": {
										"type": "keyword"
									},
									"screenshots": {
										"type": "keyword",
										"index": false
									},
									"videos": {
										"type": "keyword",
										"index": false
									},
									"awards": {
										"type": "keyword"
									}
								}
							}
						}
					}
				],
				"properties": {
					"fileRawURL": {
						"type": "text",
						"index": false
					},
					"it-riuso-codiceIPA-label": {
						"type": "text"
					},
					"publiccode-yaml-version": {
						"type": "text",
						"index": false
					},
					"name": {
						"type": "text"
					},
					"applicationSuite": {
						"type": "text",
						"fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
					},
					"url": {
						"type": "text",
						"index": false,
						"fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
					},
					"landingURL": {
						"type": "text",
						"index": false
					},
					"isBasedOn": {
						"type": "text",
						"index": false
					},
					"softwareVersion": {
						"type": "keyword"
					},
					"releaseDate": {
						"type": "date",
						"format": "strict_date"
					},
					"logo": {
						"type": "text",
						"index": false
					},
					"monochromeLogo": {
						"type": "text",
						"index": false
					},
					"inputTypes": {
						"type": "keyword"
					},
					"outputTypes": {
						"type": "keyword"
					},
					"platforms": {
						"type": "keyword"
					},
					"tags": {
						"type": "keyword"
					},
					"usedBy": {
						"type": "text",
						"fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
					},
					"roadmap": {
						"type": "text",
						"index": false
					},
					"developmentStatus": {
						"type": "keyword"
					},
					"softwareType": {
						"type": "keyword"
					},
					"intendedAudience-onlyFor": {
						"type": "keyword"
					},
					"intendedAudience-countries": {
						"type": "keyword"
					},
					"intendedAudience-unsupportedCountries": {
						"type": "keyword"
					},
					"legal-license": {
						"type": "text",
						"fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
					},
					"legal-mainCopyrightOwner": {
						"type": "text"
					},
					"legal-repoOwner": {
						"type": "text"
					},
					"legal-authorsFile": {
						"type": "text",
						"index": false
					},
					"maintenance-type": {
						"type": "keyword"
					},
					"maintenance-contractors": {
						"type": "nested",
						"properties": {
							"name": {
								"type": "text"
							},
							"until": {
								"type": "date",
								"format": "strict_date"
							},
							"website": {
								"type": "text",
								"index": false
							}
						}
					},
					"maintenance-contacts": {
						"type": "nested",
						"properties": {
							"name": {
								"type": "text"
							},
							"email": {
								"type": "text"
							},
							"phone": {
								"type": "text",
								"index": false
							},
							"affiliation": {
								"type": "text"
							}
						}
					},
					"localisation-localisationReady": {
						"type": "boolean"
					},
					"localisation-availableLanguages": {
						"type": "keyword"
					},
					"dependsOn-open": {
						"type": "nested",
						"properties": {
							"name": {
								"type": "text"
							},
							"version-min": {
								"type": "text",
								"index": false
							},
							"version-max": {
								"type": "text",
								"index": false
							},
							"version": {
								"type": "text",
								"index": false
							},
							"optional": {
								"type": "boolean"
							}
						}
					},
					"dependsOn-proprietary": {
						"type": "nested",
						"properties": {
							"name": {
								"type": "text"
							},
							"version-min": {
								"type": "text",
								"index": false
							},
							"version-max": {
								"type": "text",
								"index": false
							},
							"version": {
								"type": "text",
								"index": false
							},
							"optional": {
								"type": "boolean"
							}
						}
					},
					"dependsOn-hardware": {
						"type": "nested",
						"properties": {
							"name": {
								"type": "text"
							},
							"version-min": {
								"type": "text",
								"index": false
							},
							"version-max": {
								"type": "text",
								"index": false
							},
							"version": {
								"type": "text",
								"index": false
							},
							"optional": {
								"type": "boolean"
							}
						}
					},
					"it-conforme-accessibile": {
						"type": "boolean"
					},
					"it-conforme-interoperabile": {
						"type": "boolean"
					},
					"it-conforme-sicuro": {
						"type": "boolean"
					},
					"it-conforme-privacy": {
						"type": "boolean"
					},
					"it-spid": {
						"type": "boolean"
					},
					"it-cie": {
						"type": "boolean"
					},
					"it-anpr": {
						"type": "boolean"
					},
					"it-pagopa": {
						"type": "boolean"
					},
					"it-riuso-codiceIPA": {
						"type": "keyword"
					},
					"it-ecosistemi": {
						"type": "keyword"
					},
					"it-designKit-seo": {
						"type": "boolean"
					},
					"it-designKit-ui": {
						"type": "boolean"
					},
					"it-designKit-web": {
						"type": "boolean"
					},
					"it-designKit-content": {
						"type": "boolean"
					},
					"suggest-name": {
						"type": "completion"
					},
					"vitality-score": {
						"type": "text",
						"index": false
					},
					"vitality-dataChart": {
						"type": "integer"
					},
					"related-software": {
						"properties": {
							"name": {
								"type": "text",
								"index": false
							},
							"image": {
								"type": "text",
								"index": false
							},
							"eng": {
								"properties": {
									"localised-name": {
										"type": "text",
										"index": false
									},
									"url": {
										"type": "text",
										"index": false
									}
								}
							},
							"ita": {
								"properties": {
									"localised-name": {
										"type": "text",
										"index": false
									},
									"url": {
										"type": "text",
										"index": false
									}
								}
							}
						}
					},
					"tags-related": {
						"type": "keyword"
					},
					"popular-tags": {
						"type": "keyword"
					},
					"share-tags": {
						"type": "keyword"
					},
					"old-variant": {
						"properties": {
							"name": {
								"type": "text",
								"index": false
							},
							"eng": {
								"properties": {
									"localised-name": {
										"type": "text",
										"index": false
									},
									"url": {
										"type": "text",
										"index": false
									},
									"feature-list": {
										"type": "keyword",
										"index": false
									},
									"vitality-score": {
										"type": "integer",
										"index": false
									},
									"legal-repo-owner": {
										"type": "text",
										"index": false
									}
								}
							},
							"ita": {
								"properties": {
									"localised-name": {
										"type": "text",
										"index": false
									},
									"url": {
										"type": "text",
										"index": false
									},
									"feature-list": {
										"type": "keyword",
										"index": false
									},
									"vitality-score": {
										"type": "integer",
										"index": false
									},
									"legal-repo-owner": {
										"type": "text",
										"index": false
									}
								}
							}
						}
					},
					"old-feature-list": {
						"properties": {
							"ita": {
								"type": "keyword",
								"index": false
							},
							"eng": {
								"type": "keyword",
								"index": false
							}
						}
					}
				}
			}
		}
	}
	`
	)

	// Generating index with mapping.
	// Use the IndexExists service to check if a specified index exists.
	exists, err := elasticClient.IndexExists(index).Do(context.Background())
	if err != nil {
		return errors.New("cannot check if ES index exists for '" + index + "' exists: " + err.Error())
	}
	if !exists {
		_, err := elasticClient.CreateIndex(index).Body(mapping).Do(context.Background())
		if err != nil {
			return errors.New("cannot create ES index for '" + index + "': " + err.Error())
		}
	}

	return err
}

// ElasticFlush wrap the ElasticSearch flush command.
func ElasticFlush(index string, elasticClient *elastic.Client) error {
	// Flush to make sure the documents got written.
	_, err := elasticClient.Flush().Index(index).Do(context.Background())
	return err
}

// ElasticAliasUpdate update the Alias to the index.
func ElasticAliasUpdate(index, alias string, elasticClient *elastic.Client) error {
	// Retrieve all the aliases.
	res, err := elasticClient.Aliases().Index("_all").Do(context.Background())
	if err != nil {
		return err
	}
	// Range over all the aliases services.
	aliasService := elasticClient.Alias()
	indices := res.IndicesByAlias(alias)
	for _, name := range indices {
		log.Debugf("Remove alias from %s to %s", alias, name)
		// Remove the publiccode alias.
		_, err := aliasService.Remove(name, alias).Do(context.Background())
		if err != nil {
			return err
		}
	}

	// Add an alias to the new index.
	log.Debugf("Add alias from %s to %s", index, alias)
	_, err = aliasService.Add(index, alias).Do(context.Background())

	return err
}

// ElasticRetrier implements the elastic interface that user can implement to intercept failed requests.
type ElasticRetrier struct {
	backoff elastic.Backoff
}

// NewESRetrier returns a new ElasticRetrier with Exponential Backoff waiting.
func NewESRetrier() *ElasticRetrier {
	return &ElasticRetrier{
		backoff: elastic.NewExponentialBackoff(10*time.Millisecond, 8*time.Second),
	}
}

// Retry is used in ElasticRetrier and returns the time to wait and if the retries should stop.
func (r *ElasticRetrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	log.Warn("Elasticsearch connection problem. Retry.")

	// Stop after 8 retries: ~2m.
	if retry >= 8 {
		return 0, false, errors.New("elasticsearch or network down")
	}

	// Let the backoff strategy decide how long to wait and whether to stop.
	wait, stop := r.backoff.Next(retry)
	return wait, stop, nil
}
