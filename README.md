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
- [Træfik](https://traefik.io) as a reverse proxy

### Dependencies

- [Docker](https://www.docker.com/)
- [Docker-compose](https://docs.docker.com/compose/)
- [Go](https://golang.org/) >= 1.11

### Set-up

#### Stack

1. set up Træfik

    If you already have a Træfik container running on your host simply remove the `proxy` definition from
    `docker-compose.yml` file and set up the `web` network to be external:

    ```yaml
    networks:
      web:
        external:
          name: name_of_træfik_network
    ```

2. rename .env.example to .env and fill the variables with your values

    - default Elasticsearch user and password are `elastic:elastic`
    - default Kibana user and password are `kibana:kibana`
    - basic authentication token is generated with: `echo -n "user:password" | openssl base64 -base64`

3. rename `elasticsearch/config/searchguard/sg_internal_users.yml.example` to `elasticsearch/config/searchguard/sg_internal_users.yml` and insert the correct passwords

    Hashed passwords can be generated with:

    ```bash
    docker exec -t -i developers-italia-backend_elasticsearch elasticsearch/plugins/search-guard-6/tools/hash.sh -p <password>
    ```

4. rename config.toml.example to config.toml and fill the variables with your values

5. add mapping in `/etc/hosts` for exposed services

    For example, if `DOMAIN` in `.env` is `developers.loc`, add (if your Docker daemon is listening on localhost):

    ```
    127.0.0.1 elasticsearch.developers.loc
    127.0.0.1 kibana.developers.loc
    127.0.0.1 prometheus.developers.loc
    ```

    Or use a local DNS (like [dnsmasq](https://en.wikipedia.org/wiki/Dnsmasq)) to resolve all DNS request to `.loc` domains to localhost.

6. start the Docker stack: `make up`

#### Crawler

1. `cd crawler`
2. Fill your domains.yml file with configuration values (like specific host basic auth token)
3. Rename config.toml.example to config.toml and fill the variables with your values

##### With docker-compose (for production)

* build the crawler image: `make build`
* rename docker-compose-crawler.yml.example to docker-compose-crawler.yml. Setup the volumes mapping. Replace `network_created_by_docker_compose_prod` with the correct network name
* run `make crawl` (and configure it in crontab)

##### As golang binary (for development)

* build the crawler binary: `go build -o bin/crawler`
* start the crawler: `bin/crawler crawl whitelistPA.yml whitelistGeneric.yml`

## Troubleshooting

- From docker logs seems that Elasticsearch container needs more virtual memory and now it's `Stalling for Elasticsearch....`

  Increase container virtual memory: https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-prod-mode

- When trying to `make build` the crawler image, a fatal memory error occurs: "fatal error: out of memory"

  Probably you should increase the container memory:
  `docker-machine stop && VBoxManage modifyvm default --cpus 2 && VBoxManage modifyvm default --memory 2048 && docker-machine stop`

## Authors

[Developers Italia](https://developers.italia.it) is a project by [AgID](https://www.agid.gov.it/) in collaboration with the [Italian Digital Team](https://teamdigitale.governo.it/), which maintains this repository.