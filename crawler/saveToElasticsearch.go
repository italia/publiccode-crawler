package crawler

import (
	"context"

	"github.com/italia/developers-italia-backend/jekyll"
	"github.com/italia/developers-italia-backend/metrics"
	"github.com/olivere/elastic"
	yaml "gopkg.in/yaml.v1"
)

// File is a generic structure for saveToES() function.
// TODO: Will be replaced with a parsed publiccode.PublicCode whit proper mapping.
type File struct {
	Source string `json:"source"`
	Name   string `json:"name"`
	Data   string `json:"data"`
}

// SaveToES save the chosen data []byte in elasticsearch
func SaveToES(domain Domain, name string, data []byte, index string, elasticClient *elastic.Client) error {
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

	// Starting with elastic.v5, you must pass a context to execute each service.
	ctx := context.Background()

	// Generic publiccode data
	pc := PublicCode{}
	err := yaml.Unmarshal([]byte(data), &pc)
	if err != nil {
		return err
	}

	// Add a document to the index.
	file := jekyll.PublicCode{
		PubliccodeYamlVersion: pc.PubliccodeYamlVersion,
		Name: pc.Name,
	} //File{Source: domain.Host, Name: name, Data: string(data)}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := elasticClient.IndexExists(index).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		_, err := elasticClient.CreateIndex(index).Body(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			return err
		}
	}

	// Put publiccode data in ES.
	_, err = elasticClient.Index().
		Index(index).
		Type("software").
		Id(domain.Host + "/" + name + "_" + index).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		return err
	}

	metrics.GetCounter("repository_file_indexed", index).Inc()

	return nil
}
