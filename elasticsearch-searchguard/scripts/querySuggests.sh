#!/bin/bash
#
# Query for suggestion.
#

source config.sh

INDEX=$1

if [ ! -n "${INDEX}" ] ; then
    echo -e $RED "You have to pass index name as first parameter of the script" $Z;
    exit 1;
fi

curl -u "$BASICAUTH" -X POST "$ELASTICSEARCH_URL/$INDEX/_search?pretty" -H 'Content-Type: application/json' -d'
{
  "suggest": {
      "names" : {
          "prefix" : "Med", 
          "completion" : { 
            "field" : "suggest-name",
            "size": 10
          }
      }
  }
}
'
