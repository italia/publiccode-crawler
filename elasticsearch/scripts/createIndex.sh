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
      "properties": {
        "standard version": {
          "type": "keyword"
        },
        "url": {
          "type": "keyword",
          "index": false
        },
        "upstream-url": {
          "type": "keyword",
          "index": false
        },
        "license": {
          "type": "text",
          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
        },
        "main-copyright-owner": {
          "type": "text"
        },
        "authors-file": {
          "type": "keyword",
          "index": false
        },
        "repo-owner": {
          "type": "text"
        },
        "maintainance-type": {
          "type": "keyword"
        },
        "maintainance-until": { 
          "type": "date",
          "format": "strict_date"
        },
        "technical-contacts": {
          "properties": {
            "affiliation": {
              "type": "text",
              "fields": {
                "keyword": { "type": "keyword", "ignore_above": 256 }
              }
            },
            "email": {
              "type": "keyword"
            },
            "name": {
              "type": "text"
            }
          }
        },
        "name": {
          "type": "text"
        },
        "logo": {
          "type": "keyword",
          "index": false
        },
        "version": {
          "type": "keyword"
        },
        "released": {
          "type": "date",
          "format": "strict_date"
        },
        "platforms": {
          "type": "keyword"
        },
        "longdesc-it": {
          "type": "text"
        },
        "longdesc-en": {
          "type": "text"
        },
        "shortdesc-it": {
          "type": "text"
        },
        "shortdesc-en": {
          "type": "text"
        },
        "videos": {
          "type": "keyword",
          "index": false
        },
        "scope": {
          "type": "keyword"
        },
        "pa-type": {
          "type": "keyword"
        },
        "category": {
          "type": "keyword"
        },
        "tags": {
          "type": "keyword"
        },
        "used-by": {
          "type": "text",
          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
        },
        "dependencies": {
          "type": "text",
          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
        },
        "dependencies-hardware": {
          "type": "text",
          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
        },
        "maintainance-maintainer": {
          "type": "text",
          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
        },
        "it-use-spid": {
          "type": "keyword"
        },
        "it-use-pagopa": {
          "type": "keyword"
        },
        "suggest-name": {
          "type": "completion"
        }
      }
    }
  }
}
EOF
}

curl -u "$BASICAUTH" -X PUT "$ELASTICSEARCH_URL/$INDEX" -H 'Content-Type: application/json' -d"$(generate_index_settings)"

