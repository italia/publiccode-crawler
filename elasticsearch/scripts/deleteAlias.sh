#!/bin/bash
#
# To delete an alias from elasticsearch
#

source config.sh

ALIAS=$1
INDEX=$2

if [ ! -n "${ALIAS}" ] ; then
    echo -e $RED "You have to pass alias name as first parameter of the script" $Z;
    exit 1;
fi
if [ ! -n "${INDEX}" ] ; then
    echo -e $RED "You have to pass index name as second parameter of the script" $Z;
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

curl -u "$BASICAUTH" -X POST "http://elasticsearch:9200/_aliases" -H 'Content-Type: application/json' -d"$(generate_delete_msg)"