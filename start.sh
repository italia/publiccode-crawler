#!/usr/bin/env bash

bin/crawler updateipa
bin/crawler download-whitelist https://onboarding.developers.italia.it/repo-list whitelist/00-onboarding-reuse.yml
bin/crawler crawl whitelist/*.yml
