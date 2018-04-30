# developers-italia-backend

[![CircleCI](https://circleci.com/gh/italia/developers-italia-backend/tree/master.svg?style=svg)](https://circleci.com/gh/italia/developers-italia-backend/tree/master)

Backend &amp; crawler for the OSS catalog of Developers Italia.

**This document is a Work in progress!**

## Components

* Elasticsearch
* Kibana
* Prometheus
* Træfik

## How to contribute

### Dependencies

* [Go](https://golang.org/)
* [dep](https://github.com/golang/dep)
* [Docker](https://www.docker.com/)
* [Docker-compose](https://docs.docker.com/compose/)

### Setup

#### Stack

##### 1) rename docker-compose.yml.dist to docker-compose.yml

##### 2) set up Træfik
If you already have a Træfik container running on your host simply remove the `proxy` definition from
`docker-compose.yml` file and set up the `web` network to be external:

```yaml
networks:
  web:
    external:
      name: nome_of_træfik_network
```

##### 3) rename .env.example to .env and fill the variables with your values

* default Elasticsearch user and password are `elastic`
* default Kibana user and password are `kibana`
* basic authentication token is generated with: `echo -n "user:password" | openssl base64 -base64`

##### 4) add mapping in `/etc/hosts` for exposed services

For example, if `PROJECT_BASE_URL` in `.env` is `developers.loc`, add (if your Docker daemon is listening on localhost):
* 127.0.0.1    elasticsearch.developers.loc
* 127.0.0.1    kibana.developers.loc
* 127.0.0.1    prometheus.developers.loc

Or use a local DNS (like [dnsmasq](https://en.wikipedia.org/wiki/Dnsmasq)) to resolve all DNS request to `.loc` domains
to localhost.

##### 5) start the Docker stack: `make up`

#### Crawler

* Fill your domains.yml file with configuration values (like specific host basic auth token)

##### With docker-compose

* build the crawler image: `make build`
* de-comment the crawler container from docker-compose.yml file
* start the Docker stack: `make up`

##### As golang app

* start the crawler: `make run-all`

### Copyright

```
Copyright (c) the respective contributors, as shown by the AUTHORS file.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
