#!/bin/bash
#
# To create a index in elasticsearch
#

# shards & replicas are default values.
#

INDEX=$1

if [ ! -n "${INDEX}" ] ; then
    echo -e $RED "Devi passarmi il nome dell'indice" $Z;
    exit 1;
fi

curl -u elastic:elastic -X DELETE "http://elasticsearch:9200/$INDEX"
