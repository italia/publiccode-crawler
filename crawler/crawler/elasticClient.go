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
		elastic.SetHealthcheck(false),
	)
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
  "settings": {
    "analysis": {
      "analyzer": {
        "autocomplete": {
          "tokenizer": "autocomplete",
          "filter": [
            "lowercase"
          ]
        },
        "autocomplete_search": {
          "tokenizer": "lowercase"
        }
      },
      "tokenizer": {
        "autocomplete": {
          "type": "edge_ngram",
          "min_gram": 3,
          "max_gram": 30,
          "token_chars": [
            "letter"
          ]
        }
      }
    }
  },
  "mappings": {
    "software": {
      "dynamic_templates": [
        {
          "description": {
            "path_match": "publiccode.description.*",
            "mapping": {
              "type": "object",
              "properties": {
                "localisedName": {
                  "type": "text",
                  "analyzer": "autocomplete",
                  "search_analyzer": "autocomplete_search"
                },
                "genericName": {
                  "type": "text",
                  "fields": {
                    "keyword": { "type": "keyword", "ignore_above": 256 }
                  }
                },
                "shortDescription": {
                  "type": "text",
                  "analyzer": "autocomplete",
                  "search_analyzer": "autocomplete_search"
                },
                "longDescription": {
                  "type": "text",
                  "analyzer": "autocomplete",
                  "search_analyzer": "autocomplete_search"
                },
                "documentation": {
                  "type": "text",
                  "index": false
                },
                "apiDocumentation": {
                  "type": "text",
                  "index": false
                },
                "features": {
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
        },
        {
          "dynamic-suggestions": {
            "match": "suggest-*",
            "mapping": {
              "type": "completion",
              "preserve_separators": false
            }
          }
        }
      ],
      "properties": {
        "fileRawURL": {
          "type": "keyword",
          "index": true
        },
				"id": {
					"type": "keyword",
					"index": true
				},
				"crawltime": {
					"type": "date",
					"index": false
				},
        "publiccodeYmlVersion": {
          "type": "keyword",
          "index": false
        },

        "publiccode": {
          "properties": {
            "name": {
              "type": "text",
              "analyzer": "autocomplete",
              "search_analyzer": "autocomplete_search"
            },
            "applicationSuite": {
              "type": "text",
              "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
            },
            "url": {
              "type": "keyword",
              "index": true,
              "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
            },
            "landingURL": {
              "type": "keyword",
              "index": false
            },
            "isBasedOn": {
              "type": "keyword",
              "index": true
            },
            "softwareVersion": {
              "type": "keyword"
            },
            "releaseDate": {
              "type": "date",
              "format": "strict_date"
            },
            "logo": {
              "type": "keyword",
              "index": false
            },
            "monochromeLogo": {
              "type": "keyword",
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
              "type": "keyword",
              "index": false
            },
            "developmentStatus": {
              "type": "keyword"
            },
            "softwareType": {
              "type": "keyword"
            },
            "intendedAudience": {
              "properties": {
                "onlyFor": {
                  "type": "keyword"
                },
                "countries": {
                  "type": "keyword"
                },
                "unsupportedCountries": {
                  "type": "keyword"
                }
              }
            },
            "legal": {
              "properties": {
                "license": {
                  "type": "keyword",
                  "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
                },
                "mainCopyrightOwner": {
                  "type": "keyword"
                },
                "repoOwner": {
                  "type": "keyword"
                },
                "authorsFile": {
                  "type": "keyword",
                  "index": false
                }
              }
            },
            "maintainance": {
              "properties": {
                "type": {
                  "type": "keyword"
                },
                "contractors": {
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
                "contacts": {
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
                }
              }
            },
            "localisation": {
              "properties": {
                "localisationReady": {
                  "type": "boolean"
                },
                "availableLanguages": {
                  "type": "keyword"
                }
              }
            },
            "dependsOn": {
              "properties": {
                "open": {
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
                "proprietary": {
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
                "hardware": {
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
                }
              }
            },
            "it": {
              "properties": {
                "conforme": {
                  "properties": {
                    "accessibile": {
                      "type": "boolean"
                    },
                    "interoperabile": {
                      "type": "boolean"
                    },
                    "sicuro": {
                      "type": "boolean"
                    },
                    "privacy": {
                      "type": "boolean"
                    }
                  }
                },
                "spid": {
                  "type": "boolean"
                },
                "cie": {
                  "type": "boolean"
                },
                "anpr": {
                  "type": "boolean"
                },
                "pagopa": {
                  "type": "boolean"
                },
                "riuso": {
                  "properties": {
                    "codiceIPA": {
                      "type": "keyword"
                    }
                  }
                },
                "ecosistemi": {
                  "type": "keyword"
                },
                "designKit": {
                  "properties": {
                    "seo": {
                      "type": "boolean"
                    },
                    "ui": {
                      "type": "boolean"
                    },
                    "web": {
                      "type": "boolean"
                    },
                    "content": {
                      "type": "boolean"
                    }
                  }
                }
              }
            }
          }
        },

        "suggest-name": {
          "type": "completion"
        },
        "vitalityScore": {
          "type": "integer"
        },
        "vitalityDataChart": {
          "type": "integer"
        }
      }
    }
  }
}`
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

// ElasticAdministrationsMapping adds (if not exists) the mapping for the whitelist administrations in ES.
func ElasticAdministrationsMapping(index string, elasticClient *elastic.Client) error {
	const (
		// Elasticsearch mapping for administrations.
		mapping = `{
  "settings": {
    "analysis": {
      "analyzer": {
        "autocomplete": {
          "tokenizer": "autocomplete",
          "filter": [
            "lowercase"
          ]
        },
        "autocomplete_search": {
          "tokenizer": "lowercase"
        }
      },
      "tokenizer": {
        "autocomplete": {
          "type": "edge_ngram",
          "min_gram": 3,
          "max_gram": 30,
          "token_chars": [
            "letter"
          ]
        }
      }
    }
  },
  "mappings": {
    "administration": {
      "properties": {
        "it-riuso-codiceIPA": {
          "type": "keyword"
        },
        "it-riuso-codiceIPA-label": {
          "type": "text",
          "analyzer": "autocomplete",
          "search_analyzer": "autocomplete_search"
        }
      }
    }
  }
}`
	)

	// Generating index with mapping.
	// Use the IndexExists service to check if a specified index exists.
	exists, err := elasticClient.IndexExists(index).Do(context.Background())
	if err != nil {
		return errors.New("cannot check if ES index exists for '" + index + "' exists: " + err.Error())
	}
	// Generate new index.
	if !exists {
		_, err = elasticClient.CreateIndex(index).Body(mapping).Do(context.Background())
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
	log.Errorf("Alias ElasticAliasUpdate index-alias: %v - %v", index, alias)
	// Range over all the indices for alias service.
	aliasService := elasticClient.Alias()
	// Add an alias to the new index.
	log.Debugf("Add alias from %s to %s", index, alias)
	_, err := aliasService.Add(index, alias).Do(context.Background())

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
