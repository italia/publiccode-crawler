#!/usr/bin/env bash

# Sometimes, Elasticsearch is still not ready, yet.
# Let's wait 30 seconds to be sure Elasticsearch is ready
# to accept connections
time=30

echo "${0##*/}: Waiting ${time} seconds before running the crawler..."

sleep ${time}

bin/crawler updateipa
bin/crawler download-whitelist https://onboarding.developers.italia.it/repo-list whitelist/00-onboarding-reuse.yml
bin/crawler crawl whitelist/*.yml
