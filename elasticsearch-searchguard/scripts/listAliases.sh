#!/bin/bash
#
# Lists aliases on elasticsearch
#

source config.sh

curl -u "$BASICAUTH" -X GET "$ELASTICSEARCH_URL/_cat/aliases?v"
