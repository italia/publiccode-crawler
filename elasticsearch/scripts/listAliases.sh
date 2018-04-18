#!/bin/bash
#
# Lists aliases on elasticsearch
#

source config.sh

curl -u "$BASICAUTH" -X GET "http://elasticsearch:9200/_cat/aliases?v"
