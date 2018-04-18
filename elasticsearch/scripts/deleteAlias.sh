#!/bin/bash
#
# To create a index in elasticsearch
#

# shards & replicas are default values.
#

ALIAS=$1
INDEX=$2

if [ ! -n "${ALIAS}" ] ; then
    echo -e $RED "Devi passarmi il nome dell'Alias" $Z;
    exit 1;
fi
if [ ! -n "${INDEX}" ] ; then
    echo -e $RED "Devi passarmi il nome dell'Indice" $Z;
    exit 1;
fi

generate_delete_msg() {
  cat <<EOF
{
    "actions" : [
      { "remove" : { "index" : "$INDEX", "alias" : "$ALIAS" } }
    ]
}
EOF
}

curl -u elastic:elastic -X POST "http://elasticsearch:9200/_aliases" -H 'Content-Type: application/json' -d"$(generate_delete_msg)"