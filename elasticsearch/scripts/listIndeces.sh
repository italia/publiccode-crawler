#!/bin/bash
#
# Lists indexes in elasticsearch
#

source config.sh

curl -u "$BASICAUTH" -X GET "$ELASTICSEARCH_URL/_cat/indices?v"
