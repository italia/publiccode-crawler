#!/bin/bash
#
# Lists indexes in elasticsearch
#

source config.sh

curl -u "$BASICAUTH" -X GET "http://elasticsearch:9200/_cat/indices?v"
