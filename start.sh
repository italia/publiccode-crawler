#!/usr/bin/env bash

bin/crawler updateipa
bin/crawler download-publishers https://onboarding.developers.italia.it/repo-list publishers.onboarding.yml
bin/crawler crawl publishers*.yml
