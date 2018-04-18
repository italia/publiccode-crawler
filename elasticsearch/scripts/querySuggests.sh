#!/bin/bash
#
# Query for suggestion.
#

INDEX=$1

if [ ! -n "${INDEX}" ] ; then
    echo -e $RED "Devi passarmi il nome dell'Indice" $Z;
    exit 1;
fi

curl -u elastic:elastic -X POST "http://elasticsearch:9200/$INDEX/_search?pretty" -H 'Content-Type: application/json' -d'
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
