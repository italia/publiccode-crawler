# publiccode.yml crawler for the software catalog of Developers Italia

[![Go Report Card](https://goreportcard.com/badge/github.com/italia/publiccode-crawler/v3)](https://goreportcard.com/report/github.com/italia/publiccode-crawler/v3)
[![Join the #publiccode channel](https://img.shields.io/badge/Slack%20channel-%23publiccode-blue.svg?logo=slack)](https://developersitalia.slack.com/messages/CAM3F785T)
[![Get invited](https://slack.developers.italia.it/badge.svg)](https://slack.developers.italia.it/)

## Description

Developers Italia provides [a catalog of Free and Open Source](https://developers.italia.it/en/search)
software aimed to Public Administrations.

This **crawler** retrieves the `publiccode.yml` files from the
repositories of publishers found in the [Developers Italia API](https://github.com/italia/developers-italia-api).

## Setup and deployment processes

The crawler can either run manually on the target machine or it can be deployed
from a Docker container.

### Manually configure and build the crawler

1. Rename `config.toml.example` to `config.toml` and set the variables

   > **NOTE**: The application also supports environment variables in substitution
   > to config.toml file. Remember: "environment variables get higher priority than
   > the ones in configuration file"

2. Build the binary with `go build`

### Docker

You can build the Docker image using

```console
docker build .
```

or use the image published to DockerHub:

```console
docker run -it italia/publiccode-crawler
```

## Commands

### `crawler crawl`

Gets the list of publishers from `https://api.developers.italia.it/v1/publishers`
and starts to crawl their repositories.

### `crawler crawl publishers*.yml`

Gets the list of publishers in `publishers*.yml` and starts to crawl
their repositories.

### Other commands

* `crawler download-publishers` downloads organizations and repositories from
  the [onboarding portal repository](https://github.com/italia/developers-italia-onboarding)
  and saves them to a publishers YAML file.

## See also

* [developers-italia-api](https://github.com/italia/developers-italia-api): the API
  used to store the results of the crawling
* [publiccode-parser-go](https://github.com/italia/publiccode-parser-go): the Go
  package for parsing publiccode.yml files

## Authors

[Developers Italia](https://developers.italia.it) is a project by
[AgID](https://www.agid.gov.it/) and the
[Italian Digital Team](https://teamdigitale.governo.it/), which developed the
crawler and maintains this repository.
