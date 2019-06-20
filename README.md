# developers-italia-backend

[![CircleCI](https://circleci.com/gh/italia/developers-italia-backend/tree/master.svg?style=shield)](https://circleci.com/gh/italia/developers-italia-backend/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/italia/developers-italia-backend)](https://goreportcard.com/report/github.com/italia/developers-italia-backend) [![Join the #website channel](https://img.shields.io/badge/Slack%20channel-%23website-blue.svg?logo=slack)](https://developersitalia.slack.com/messages/C9R26QMT6)
[![Get invited](https://slack.developers.italia.it/badge.svg)](https://slack.developers.italia.it/)

## Backend & crawler for the OSS catalog of Developers Italia

The crawler finds and retrieves all the publiccode.yml files from the Organizations registered on Github/Bitbucket/Gitlab listed in the whitelistes, and then generates YAML files that are later used by the [Jekyll build chain](https://github.com/italia/developers.italia.it) to generate the static pages of [developers.italia.it](https://developers.italia.it/).

### Components

- [Elasticsearch](https://www.elastic.co/products/elasticsearch) for storing the data
- [Kibana](https://www.elastic.co/products/kibana) for internal visualization of data
- [Prometheus](https://prometheus.io) for collecting metrics

### Dependencies

- [Docker](https://www.docker.com/)
- [Docker-compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/) >= 1.11

### Set-up

#### Stack

1. rename .env.example to .env and fill the variables with your values

    - default Elasticsearch user and password are `elastic:elastic`
    - default Kibana user and password are `kibana:kibana`

2. rename `elasticsearch/config/searchguard/sg_internal_users.yml.example` to `elasticsearch/config/searchguard/sg_internal_users.yml` and insert the correct passwords

    Hashed passwords can be generated with:

    ```bash
    docker exec -t -i developers-italia-backend_elasticsearch elasticsearch/plugins/search-guard-6/tools/hash.sh -p <password>
    ```

3. insert the `kibana` password in `kibana/config/kibana.yml`

4. configure the nginx proxy for the lasticsearch host with the following directives:

    ```
    limit_req_zone $binary_remote_addr zone=elasticsearch_limit:10m rate=10r/s;

    server {
        ...
        location / {
            limit_req zone=elasticsearch_limit burst=20 nodelay;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_pass http://localhost:9200;
            proxy_ssl_session_reuse off;
            proxy_cache_bypass $http_upgrade;
            proxy_redirect off;
        }
    }
    ```

5. you might need to type `sysctl -w vm.max_map_count=262144` and make this permanent in /etc/sysctl.conf in order to start elasticsearch, as [documented here](https://hub.docker.com/r/khezen/elasticsearch/)

6. start the Docker stack: `make up`

#### Crawler

1. `cd crawler`
2. Fill your domains.yml file with configuration values (like specific host basic auth tokens)
3. Rename config.toml.example to config.toml and fill the variables
4. build the crawler binary: `make`
5. start the crawler: `bin/crawler crawl whitelist/*.yml`
6. configure in crontab as desired

### Tools

* `bin/crawler updateipa` downloads IPA data and writes it into Elasticsearch
* `bin/crawler download-whitelist` downloads orgs and repos from the [onboarding portal](https://github.com/italia/developers-italia-onboarding) and writes them to a whitelist file

### Troubleshooting

- From docker logs seems that Elasticsearch container needs more virtual memory and now it's `Stalling for Elasticsearch....`

  Increase container virtual memory: https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-prod-mode

- When trying to `make build` the crawler image, a fatal memory error occurs: "fatal error: out of memory"

  Probably you should increase the container memory:
  `docker-machine stop && VBoxManage modifyvm default --cpus 2 && VBoxManage modifyvm default --memory 2048 && docker-machine stop`

### Development

In order to access Elasticsearch with write permissions from the outside, you can forward the 9200 port via SSH using `ssh -L9200:localhost:9200` and configure `ELASTIC_URL = "http://localhost:9200/"` in your local config.toml.

## See also

* [publiccode-parser-go](https://github.com/italia/publiccode-parser-go): the Go package for parsing publiccode.yml files
* [developers-italia-onboarding](https://github.com/italia/developers-italia-onboarding): the onboarding portal

## Authors

[Developers Italia](https://developers.italia.it) is a project by [AgID](https://www.agid.gov.it/) and the [Italian Digital Team](https://teamdigitale.governo.it/), which developed the crawler and maintains this repository.