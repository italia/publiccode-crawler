#!/bin/bash
#
# To create an index in elasticsearch
#

# shards & replicas are default values.
#

source config.sh

TODAY=$(date '+%Y%m%d')
INDEX="publiccode_$TODAY"

generate_index_settings() {
  cat <<EOF
{
  "settings" : {
    "index" : {
      "number_of_shards" : 5,
      "number_of_replicas" : 1
    }
  },
  "mappings": {
    "software": {
      "dynamic_templates": [
        {
          "free-tags": {
            "match_pattern": "regex",
            "match": "^free-tags-[a-z]{3}$",
            "mapping": {
              "type": "keyword"
            }
          }
        },
        {
          "description": {
            "path_match": "description.*",
            "mapping": {
              "type": "object",
              "properties": {
                "localised-name": {
                  "type": "text"
                },
                "short-description": {
                  "type": "text"
                },
                "long-description": {
                  "type": "text"
                },
                "documentation": {
                  "type": "keyword",
                  "index": false
                },
                "feature-list": {
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
        "application-suite": {
          "type": "text",
          "fields": { "keyword":{ "type": "keyword", "ignore_above": 256 } }
        },
        "url": {
          "type": "text",
          "index": false,
          "fields": { "keyword":{ "type": "keyword", "ignore_above": 256 } }
        },
        "landing-url": {
          "type": "text",
          "index": false
        },
        "is-based-on": {
          "type": "text",
          "index": false
        },
        "software-version": {
          "type": "keyword"
        },
        "release-date": {
          "type": "date",
          "format": "strict_date"
        },
        "logo": {
          "type": "text",
          "index": false
        },
        "monochrome-logo": {
          "type": "text",
          "index": false
        },
        "platforms": {
          "type": "keyword"
        },
        "tags": {
          "type": "keyword"
        },
        "used-by": {
          "type": "text",
          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
        },
        "roadmap": {
          "type": "text",
          "index": false
        },
        "development-status": {
          "type": "keyword"
        },
        "software-type": {
          "type": "keyword"
        },
        "software-type-is-related-to": {
          "type": "text",
          "index": false          
        },
        "intended-audience-only-for": {
          "type": "keyword"
        },
        "intended-audience-countries": {
          "type": "keyword"
        },
        "intended-audience-unsupported-countries": {
          "type": "keyword"
        },
        "legal-license": {
          "type": "text",
          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
        },
        "legal-main-copyright-owner": {
          "type": "text"
        },
        "legal-repo-owner": {
          "type": "text"
        },
        "legal-authors-file": {
          "type": "text",
          "index": false
        },
        "maintainance-type": {
          "type": "keyword"
        },
        "maintainance-contractors": {
          "type": "nested",
          "properties": {
            "name": {
              "type":"text"
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
        "maintainance-contacts": {
          "type": "nested",
          "properties": {
            "name": {
              "type":"text"
            },
            "email": {
              "type": "text"
            },
            "phone": {
              "type": "text",
              "index": false
            },
            "affiliation": {
              "type":"text"
            }
          }
        },
        "localisation-localisation-ready": {
          "type": "boolean"
        },
        "localisation-available-languages": {
          "type": "keyword"
        },
        "dependencies-software": {
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
            "optional": {
              "type": "boolean"
            }
          }
        },
        "dependencies-hardware": {
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
            "optional": {
              "type": "boolean"
            }
          }
        },
        "it-accessibile": {
          "type":"boolean"
        },
        "it-spid": {
          "type":"boolean"
        },
        "it-cie": {
          "type":"boolean"
        },
        "it-anpr": {
          "type":"boolean"
        },
        "it-pagopa": {
          "type":"boolean"
        },
        "it-riuso-codice-ipa": {
          "type":"boolean"
        },
        "it-design-kit-service-design" : {
          "type":"boolean"
        },
        "it-design-kit-ui" : {
          "type":"boolean"
        },
        "it-design-kit-web-toolkit" : {
          "type":"boolean"
        },        
        "suggest-name": {
          "type": "completion"
        },
        "vitality-score": {
          "type": "text",
          "index": false
        },
        "vitality-data-chart": {
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
                }
              }
            }
          }
        }
      }
    }
  }
}
EOF
}

curl -u "$BASICAUTH" -X PUT "$ELASTICSEARCH_URL/$INDEX" -H 'Content-Type: application/json' -d"$(generate_index_settings)"

