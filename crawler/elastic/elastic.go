package elastic

import (
	"context"
	"errors"
	"net/http"
	"time"

	elastic "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ClientFactory returns an elastic Client.
func ClientFactory(URL, user, password string) (*elastic.Client, error) {
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

// PubliccodeMapping is the Elasticsearch mapping for the publiccode index.
// AdministrationsMapping is the Elasticsearch mapping for the administrations index.
const (
	PubliccodeMapping = `{
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
              "search_analyzer": "autocomplete_search",
              "fields": {
                "keyword": { "type": "keyword", "ignore_above": 256 }
              }
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
          "search_analyzer": "autocomplete_search",
          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
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
        "categories": {
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
            "scope": {
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
            "countryExtensionVersion": {
              "type": "keyword",
              "index": false
            },
            "conforme": {
              "properties": {
                "lineeGuidaDesign": {
                  "type": "boolean"
                },
                "modelloInteroperabilita": {
                  "type": "boolean"
                },
                "misureMinimeSicurezza": {
                  "type": "boolean"
                },
                "gdpr": {
                  "type": "boolean"
                }
              }
            },
            "piattaforme": {
              "properties": {
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
                }
              }
            },
            "riuso": {
              "properties": {
                "codiceIPA": {
                  "type": "keyword"
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
}`
	AdministrationsMapping = `{
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
    "properties": {
      "it-riuso-codiceIPA": {
        "type": "keyword"
      },
      "it-riuso-codiceIPA-label": {
        "type": "text",
        "analyzer": "autocomplete",
        "search_analyzer": "autocomplete_search"
      },
      "type": {
        type: "keyword"
      }
    }
  }
}`
	IPAMapping = `{
    "settings": {
      "index": {
        "analysis": {
          "filter": {},
          "analyzer": {
            "autocomplete": {
              "tokenizer": "autocomplete",
              "filter": "lowercase"
            },
            "autocomplete_search": {
              "tokenizer": "lowercase"
            }
          },
          "tokenizer": {
            "autocomplete": {
              "type": "edge_ngram",
              "min_gram": 2,
              "max_gram": 20,
              "token_chars": [
                "letter"
              ]
            }
          },
          "normalizer": {
            "lowercase_normalizer": {
              "type": "custom",
              "char_filter": [],
              "filter": ["lowercase"]
            }
          }
        }
      }
    },
    "mappings": {
      "dynamic": "strict",
      "properties": {
        "ipa": {
          "type": "text",
          "analyzer": "autocomplete",
          "search_analyzer": "autocomplete_search",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256,
              "normalizer": "lowercase_normalizer"
            }
          }
        },
        "description": {
          "type": "text",
          "analyzer": "autocomplete",
          "search_analyzer": "autocomplete_search"
        },
        "pec": {
          "type": "text",
          "analyzer": "autocomplete",
          "search_analyzer": "autocomplete_search",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256,
              "normalizer": "lowercase_normalizer"
            }
          }
        },
        "type": {
          "type": "keyword"
        },
        "cf": {
          "type": "keyword"
        },
        "website": {
          "type": "keyword"
        }
      }
    }
  }`
)

// CreateIndexMapping adds (if not exists) the mapping for the crawler data in ES.
func CreateIndexMapping(index string, mapping string, elasticClient *elastic.Client) error {
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

// Flush wrap the ElasticSearch flush command.
func Flush(index string, elasticClient *elastic.Client) error {
	// Flush to make sure the documents got written.
	_, err := elasticClient.Flush().Index(index).Do(context.Background())
	return err
}

// AliasUpdate update the Alias to the index.
func AliasUpdate(index, alias string, elasticClient *elastic.Client) error {
	// Range over all the indices for alias service.
	aliasService := elasticClient.Alias()

	// Add an alias to the new index.
	log.Debugf("Add alias from %s to %s", index, alias)
	_, err := aliasService.Add(index, alias).Do(context.Background())

	return err
}

// Retrier implements the elastic interface that user can implement to intercept failed requests.
type Retrier struct {
	backoff elastic.Backoff
}

// NewESRetrier returns a new Retrier with Exponential Backoff waiting.
func NewESRetrier() *Retrier {
	return &Retrier{
		backoff: elastic.NewExponentialBackoff(10*time.Millisecond, 8*time.Second),
	}
}

// Retry is used in Retrier and returns the time to wait and if the retries should stop.
func (r *Retrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	log.Warn("Elasticsearch connection problem. Retry.")

	// Stop after 8 retries: ~2m.
	if retry >= 8 {
		return 0, false, errors.New("elasticsearch or network down")
	}

	// Let the backoff strategy decide how long to wait and whether to stop.
	wait, stop := r.backoff.Next(retry)
	return wait, stop, nil
}

// NewBoolQuery initializes a boolean query for Elasticsearch.
func NewBoolQuery(queryType string) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	if queryType == "software" {
		unsupportedCountries := viper.GetStringSlice("IGNORE_UNSUPPORTEDCOUNTRIES")
		uc := make([]interface{}, len(unsupportedCountries))
		for i, v := range unsupportedCountries {
			uc[i] = v
		}
		query = query.MustNot(elastic.NewTermsQuery("publiccode.intendedAudience.unsupportedCountries", uc...))
	}

	return query
}
