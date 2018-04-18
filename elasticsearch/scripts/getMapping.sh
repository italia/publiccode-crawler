#!/bin/bash
#
# To get index mappings
#

INDEX=$1

if [ ! -n "${INDEX}" ] ; then
    echo -e $RED "Devi passarmi il nome dell'Indice" $Z;
    exit 1;
fi

curl -u elastic:elastic -X GET "http://elasticsearch:9200/$INDEX/_mapping"
