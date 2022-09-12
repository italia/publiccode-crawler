# publiccode.yml crawler for the software catalog of Developers Italia

[![Go Report Card](https://goreportcard.com/badge/github.com/italia/developers-italia-backend)](https://goreportcard.com/report/github.com/italia/developers-italia-backend)
[![Join the #publiccode channel](https://img.shields.io/badge/Slack%20channel-%23publiccode-blue.svg?logo=slack)](https://developersitalia.slack.com/messages/CAM3F785T)
[![Get invited](https://slack.developers.italia.it/badge.svg)](https://slack.developers.italia.it/)

## Description

Developers Italia provides [a catalog of Free and Open Source](https://developers.italia.it/en/search)
software aimed to Public Administrations.

This **crawler** retrieves the `publiccode.yml` files from the
organizations publishing the software that have registered through the
[onboarding procedure](https://github.com/italia/developers-italia-onboarding).

## Setup and deployment processes

The crawler can either run manually on the target machine or it can be deployed
from a Docker container with
[its helm-chart](https://github.com/teamdigitale/devita-infra-kubernetes) in Kubernetes.

### Manually configure and build the crawler

1. Rename `config.toml.example` to `config.toml` and set the variables

   > **NOTE**: The application also supports environment variables in substitution
   > to config.toml file. Remember: "environment variables get higher priority than
   > the ones in configuration file"

2. Build the crawler binary with `make`

### Docker

The repository has a `Dockerfile`, used to build the production image,
and a `docker-compose.yml` file to setup the development environment.

1. Start the environment:

   ```shell
   docker-compose up

## Run the crawler

### Crawl mode: `bin/crawler crawl publishers*.yml`

Gets the list of publishers in `publishers*.yml` and starts to crawl
their repositories.

### One mode (single repository url): `bin/crawler one [repo url] publishers*.yml`

In this mode one single repository at the time will be evaluated. If the
organization is present, its iPA code will be matched with the ones in
the publishers' file, otherwise it will be set to null and the `slug` will have
a random code in the end (instead of the iPA code).

Furthermore, the iPA code validation, which is a simple check within the publishers'
file (to ensure that code belongs to the selected publisher), will be skipped.

### Other commands

* `bin/crawler download-publishers` downloads organizations and repositories from
  the [onboarding portal repository](https://github.com/italia/developers-italia-onboarding)
  and saves them to a publishers YAML file.

## See also

* [publiccode-parser-go](https://github.com/italia/publiccode-parser-go): the Go
  package for parsing publiccode.yml files

## Authors

[Developers Italia](https://developers.italia.it) is a project by
[AgID](https://www.agid.gov.it/) and the
[Italian Digital Team](https://teamdigitale.governo.it/), which developed the
crawler and maintains this repository.
