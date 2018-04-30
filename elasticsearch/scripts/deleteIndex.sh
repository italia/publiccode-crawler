#!/bin/bash
#
# To delete an index from elasticsearch
#

source config.sh

INDEX=$1

if [ ! -n "${INDEX}" ] ; then
    echo -e $RED "You have to pass index name as first parameter of the script" $Z;
    exit 1;
fi

curl -u "$BASICAUTH" -X DELETE "$ELASTICSEARCH_URL/$INDEX"
