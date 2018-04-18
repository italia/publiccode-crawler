#!/bin/bash
#
# To create a index in elasticsearch
#

# shards & replicas are default values.
#

curl -u elastic:elastic -X GET "http://elasticsearch:9200/_cat/aliases?v"
