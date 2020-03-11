#!/bin/bash
#
# To get index mappings
#

source config.sh

INDEX=$1

if [ ! -n "${INDEX}" ] ; then
    echo -e $RED "You have to pass index name" $Z;
    exit 1;
fi

curl -u "$BASICAUTH" -X GET "$ELASTICSEARCH_URL/$INDEX/_mapping"
