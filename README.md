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

The generated YAML files are then used by
[developers.italia.it build](https://github.com/italia/developers.italia.it)
to generate its static pages.

## Setup and deployment processes

The crawler can either run manually on the target machine or it can be deployed
from a Docker container with
[its helm-chart](https://github.com/teamdigitale/devita-infra-kubernetes) in Kubernetes.

[Elasticsearch](https://www.elastic.co/products/elasticsearch) is used to store
the data and has ready to accept connections before the crawler is started.

### Manually configure and build the crawler

1. Save the auth tokens to `domains.yml`.

2. Rename `config.toml.example` to `config.toml` and set the variables

   > **NOTE**: The application also supports environment variables in substitution
   > to config.toml file. Remember: "environment variables get higher priority than
   > the ones in configuration file"

3. Build the crawler binary with `make`

### Docker

The repository has a `Dockerfile`, used to build the production image,
and a `docker-compose.yml` file to setup the development environment.

1. Copy the [`.env.example`](.env.example) file into `.env` and edit the
   environment variables as it suits you.
   [`.env.example`](.env.example) has detailed descriptions for each variable.

   ```shell
   cp .env.example .env
   ```

2. Save your auth tokens to `domains.yml`

   ```shell
   cp crawler/domains.yml.example crawler/domains.yml
   editor crawler/domains.yml
   ```

3. Start the environment:

   ```shell
   docker-compose up

## Run the crawler

### Crawl mode (all item in whitelists): `bin/crawler crawl whitelist/*.yml`

Gets the list of organizations in `whitelist/*.yml` and starts to crawl
their repositories.

If it finds a blacklisted repository, it will remove it from Elasticsearch, if
it is present.

It also generates:

* [`amministrazioni.yml`](https://crawler.developers.italia.it/amministrazioni.yml)
  containing all the Public Administrations their name, website URL and iPA code.

* [`softwares.yml`](https://crawler.developers.italia.it/softwares.yml) containing
  all the software that the crawler scraped, validated and saved into ElasticSearch.

  The structure is similar to publiccode data structure with some additional
  fields like vitality and vitality score.

* [`software-riuso.yml`](https://crawler.developers.italia.it/software-riuso.yml)
  containing all the software in `softwares.yml` having an iPA code.

* [`software-open-source.yml`](https://crawler.developers.italia.it/software-open-source.yml)
  containing all the software in `softwares.yml` with no iPA code.

* `https://crawler.developers.italia.it/HOSTING/ORGANIZATION/REPO/log.json` containing
  the logs of the scraping for that particular `REPO`.
  (eg. [`https://crawler.developers.italia.it/github.com/italia/design-scuole-wordpress-theme/log.json`](https://crawler.developers.italia.it/github.com/italia/design-scuole-wordpress-theme/log.json))

### One mode (single repository url): `bin/crawler one [repo url] whitelist/*.yml`

In this mode one single repository at the time will be evaluated. If the
organization is present, its iPA code will be matched with the ones in
whitelist, otherwise it will be set to null and the `slug` will have a random
code in the end (instead of the iPA code).

Furthermore, the iPA code validation, which is a simple check within whitelists
(to ensure that code belongs to the selected PA), will be skipped.

If it finds a blacklisted repository, it will exit immediately.

### Other commands

* `bin/crawler updateipa` downloads iPA data and writes them into Elasticsearch

* `bin/crawler delete [URL]` deletes software from Elasticsearch using its code
   hosting URL specified in `publiccode.url`

* `bin/crawler download-whitelist` downloads organizations and repositories from
  the [onboarding portal repository](https://github.com/italia/developers-italia-onboarding)
  and saves them to a whitelist file

### Crawler whitelists

The whitelist directory contains the of organizations to crawl from.

`whitelist/manual-reuse.yml` is a list of Public Administrations repositories
that for various reasons were not onboarded with
[developers-italia-onboarding](https://github.com/italia/developers-italia-onboarding),
while `whitelist/thirdparty.yml` contains the non-PAs repos.

Here's an example of how the files might look like:

```yaml
- id: "Comune di Bagnacavallo" # generic name of the organization.
  codice-iPA: "c_a547" # codice-iPA
  organizations: # list of organization urls.
    - "https://github.com/gith002"
```

### Crawler blacklists

Blacklists are needed to exclude individual repository that are not in line with
our
[guidelines](https://docs.italia.it/italia/developers-italia/policy-inserimento-catalogo-docs/it/stabile/approvazione-del-software-a-catalogo.html).

You can set `BLACKLIST_FOLDER` in `config.toml` to point to a directory
where blacklist files are located.
Blacklisting is currently supported by the `one` and `crawl` commands.

## See also

* [publiccode-parser-go](https://github.com/italia/publiccode-parser-go): the Go
  package for parsing publiccode.yml files

## Authors

[Developers Italia](https://developers.italia.it) is a project by
[AgID](https://www.agid.gov.it/) and the
[Italian Digital Team](https://teamdigitale.governo.it/), which developed the
crawler and maintains this repository.
