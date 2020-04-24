# Backend and crawler for the OSS catalog of Developers Italia
[![CircleCI](https://circleci.com/gh/italia/developers-italia-backend/tree/master.svg?style=shield)](https://circleci.com/gh/italia/developers-italia-backend/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/italia/developers-italia-backend)](https://goreportcard.com/report/github.com/italia/developers-italia-backend) [![Join the #website channel](https://img.shields.io/badge/Slack%20channel-%23website-blue.svg?logo=slack)](https://developersitalia.slack.com/messages/C9R26QMT6)
[![Get invited](https://slack.developers.italia.it/badge.svg)](https://slack.developers.italia.it/)

## Overview: how the crawler works

The crawler finds and retrieves the *publiccode.yml* files from the organizations registered on *Github/Bitbucket/Gitlab*, listed in the whitelist.
It then creates YAML files used by the [Jekyll build chain](https://github.com/italia/developers.italia.it) to generate the static pages of [developers.italia.it](https://developers.italia.it/).

## Dependencies and other related software

These are the dependencies and some useful tools used in conjunction with the crawler.

* [Elasticsearch 6.8.7](https://www.elastic.co/products/elasticsearch) for storing the data. Elasticsearch should be active and ready to accept connections before the crawler gets started

* [Kibana 6.8.7](https://www.elastic.co/products/kibana) for internal data visualization (optional)

* [Prometheus 6.8.7](https://prometheus.io) for collecting metrics (optional, currently supported but not used in production)

## Tools

This is the list of tools used in the repository:

* [Docker](https://www.docker.com/)

* [Docker-compose](https://docs.docker.com/compose/)

* [Go](https://golang.org/) >= 1.11

## Setup and deployment processes

The crawler can either run directly on the target machine, or it can be deployed in form of Docker container, possibly using an orchestrator, such as Kubernetes.

Up to now, the crawler and its dependencies have run in form of Docker containers on a virtual machine. Elasticsearch and Kibana have been deployed using a fork of the main project, called [search guard](https://search-guard.com/). This is still deployed in production and what we'll call in the readme *"legacy deployment process"*.

With the idea of making the legacy installation more scalable and reliable, a refactoring of the code has been recently made. The readme refers to this approach as the *new deployment process*. This includes using the official version of Elasticsearch and Kibana, and deploying the Docker containers on top of Kubernetes, using helm-charts. While the crawler has it's [own helm-chart](https://github.com/teamdigitale/devita-infra-kubernetes), Elasticsearch and Kibana are deployed using their [official helm-charts](https://github.com/elastic/helm-charts).
The new deployment process uses a [docker-compose.yml](docker-compose.yml) file to only bring up a local development and test environment.

The paragraph starts describing how to build and run the crawler, directly on a target machine.
The procedure described is the same automated in the Dockerfile. The -legacy and new- Docker deployment procedures are then described below.

### Manually configure and build the crawler

* `cd crawler`

* Fill the *domains.yml* file with configuration values (i.e. host basic auth tokens)

* Rename the *config.toml.example* file to *config.toml* and fill the variables

> **NOTE**: The application also supports environment variables in substitution to config.toml file. Remember: "environment variables get higher priority than the ones in configuration file"

* Build the crawler binary: `make`

* Configure the crontab as desired

### Run the crawler
* Crawl mode (all item in whitelists): `bin/crawler crawl whitelist/*.yml`

* One mode (single repository url): `bin/crawler one [repo url] whitelist/*.yml`
  - In this mode one single repository at the time will be evaluated. If the organization is present, its IPA code will be matched with the ones in whitelist otherwise it will be set to null and the `slug` will have a random code in the end (instead of the IPA code). Furthermore, the IPA code validation, which is a simple check within whitelists (to ensure that code belongs to the selected PA), will be skipped.

* `bin/crawler updateipa` downloads IPA data and writes them into Elasticsearch

* `bin/crawler delete [URL]` delete software from Elasticsearch using its code hosting URL specified in `publiccode.url` 

* `bin/crawler download-whitelist` downloads organizations and repositories from the [onboarding portal repository](https://github.com/italia/developers-italia-onboarding) and saves them to a whitelist file

### Docker: the legacy deployment process

The paragraph describes how to setup and deploy the crawler, following the *legacy deployment process*.

* Rename [.env-search-guard.example](.env-search-guard.example) to *.env* and adapt its variables as needed

* Rename *elasticsearch-searchguard/config/searchguard/sg_internal_users.yml.example* to *elasticsearch/-searchguard/config/searchguard/sg_internal_users.yml* and insert the correct passwords. Hashed passwords can be generated with:

    ```shell
    docker exec -t -i developers-italia-backend_elasticsearch elasticsearch-searchguard/plugins/search-guard-6/tools/hash.sh -p <password>
    ```

* Insert the *kibana* password in [kibana-searchguard/config/kibana.yml](kibana-searchguard/config/kibana.yml)

* Configure the Nginx proxy for the elasticsearch host with the following directives:

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

* You might need to type `sysctl -w vm.max_map_count=262144` and make this permanent in /etc/sysctl.conf in order to start elasticsearch, as [documented here](https://hub.docker.com/r/khezen/elasticsearch/)

* Start Docker: `make up`

### Docker: the new deployment process

The repository has a *Dockerfile*, used to also build the production image, and a *docker-compose.yml* file to facilitate the local deployment.

The containers declared in the *docker-compose.yml* file leverage some environment variables that should be declared in a *.env* file. A [.env.example](.env.example) file has some exemplar values. Before proceeding with the build, copy the [.env.example](.env.example) into *.env* and modify the environment variables as needed.

To build the crawler container, download its dependencies and start them all, run:

```shell
docker-compose up [-d] [--build]
```

where:

* *-d* execute the containers in background

* *--build* forces the containers build

To destroy the containers, use:

```shell
docker-compose down
```

#### Xpack

By default, the system -specifically Elasticsearch- doesn't make use of xpack, so passwords and certificates. To do so, the Elasticsearch container mounts [this configuration file](elasticsearch/elasticsearch.yml). This will make things work out of the box, but it's not appropriate for production environments.

An alternative configuration file that enables xpack is available [here](elasticsearch/elasticsearch-xpack.yml). In order to use it, you should

* Generate appropriate certificates for elasticsearch, save them in the *elasticsearch folder*, and make sure that their name matches the one contained in the [elasticsearch-xpack configuration file](elasticsearch/elasticsearch-xpack.yml).

* Optionally change the [elasticsearch-xpack.yml configuration file](elasticsearch/elasticsearch-xpack.yml) as desired

* Rename the [elasticsearch-xpack.yml configuration file](elasticsearch/elasticsearch-xpack.yml) to *elasticsearch.yml*

* Change the environment variables in your *.env* file to make sure that crawler, elasticsearch, and kibana configurations have matching passwords

At this point you can bring up the environment with *docker-compose*.

## Troubleshooting Q/A

* From docker logs seems that Elasticsearch container needs more virtual memory and now it's *Stalling for Elasticsearch...*

    Increase container virtual memory: https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-prod-mode

* When trying to `make build` the crawler image, a fatal memory error occurs: "fatal error: out of memory"

    Probably you should increase the container memory: `docker-machine stop && VBoxManage modifyvm default --cpus 2 && VBoxManage modifyvm default --memory 2048 && docker-machine stop`

## See also

* [publiccode-parser-go](https://github.com/italia/publiccode-parser-go): the Go package for parsing publiccode.yml files

* [developers-italia-onboarding](https://github.com/italia/developers-italia-onboarding): the onboarding portal

## Authors

[Developers Italia](https://developers.italia.it) is a project by [AgID](https://www.agid.gov.it/) and the [Italian Digital Team](https://teamdigitale.governo.it/), which developed the crawler and maintains this repository.
