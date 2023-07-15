#!/usr/bin/env bash

bin/crawler updateipa
bin/crawler download-publishers https://onboarding.developers.italia.it/repo-list publishers.onboarding.yml
wget https://raw.githubusercontent.com/italia/developers.italia.it/HEAD/_data/publishers.thirdparty.yml
wget https://raw.githubusercontent.com/italia/developers.italia.it/HEAD/_data/publishers.yml

bin/crawler crawl publishers*.yml
