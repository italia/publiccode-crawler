#!/usr/bin/env bash

curl 'https://www.indicepa.gov.it/public-services/opendata-read-service.php?dstype=FS&filename=pec.txt' > pec.txt

bin/crawler download-publishers https://onboarding.developers.italia.it/repo-list publishers.onboarding.yml
bin/crawler crawl publishers*.yml
