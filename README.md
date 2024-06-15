# publiccode.yml crawler for the software catalog of Developers Italia

[![Go Report Card](https://goreportcard.com/badge/github.com/italia/publiccode-crawler/v4)](https://goreportcard.com/report/github.com/italia/publiccode-crawler/v4)
[![Join the #publiccode channel](https://img.shields.io/badge/Slack%20channel-%23publiccode-blue.svg?logo=slack)](https://developersitalia.slack.com/messages/CAM3F785T)
[![Get invited](https://slack.developers.italia.it/badge.svg)](https://slack.developers.italia.it/)

## Description

Developers Italia provides [a catalog of Free and Open Source](https://developers.italia.it/en/search)
software aimed to Public Administrations.

`publiccode-crawler` retrieves the `publiccode.yml` files from the
repositories of publishers found in the [Developers Italia API](https://github.com/italia/developers-italia-api).

## Setup and deployment processes

`publiccode-crawler` can either run manually on the target machine or it can be deployed
from a Docker container.

### Manually configure and build

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

### `publiccode-crawler crawl`

Gets the list of publishers from `https://api.developers.italia.it/v1/publishers`
and starts to crawl their repositories.

### `publiccode-crawler crawl publishers*.yml`

Gets the list of publishers in `publishers*.yml` and starts to crawl
their repositories.

### `publiccode-crawler crawl-software <software> <publisher>`

Crawl just the software specified as parameter.
It takes the software URL and its publisher id as parameters.

Ex. `publiccode-crawler crawl-software https://api.developers.italia.it/v1/software/a2ea59b0-87cd-4419-b93f-00bed8a7b859 edb66b3d-3e36-4b69-aba9-b7c4661b3fdd`

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
