# Crawler for the OSS catalog of Developers Italia

[![CircleCI](https://circleci.com/gh/italia/developers-italia-backend/tree/master.svg?style=shield)](https://circleci.com/gh/italia/developers-italia-backend/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/italia/developers-italia-backend)](https://goreportcard.com/report/github.com/italia/developers-italia-backend)
[![Join the #website channel](https://img.shields.io/badge/Slack%20channel-%23website-blue.svg?logo=slack)](https://developersitalia.slack.com/messages/C9R26QMT6)
[![Get invited](https://slack.developers.italia.it/badge.svg)](https://slack.developers.italia.it/)

## How it works

The crawler finds and retrieves the **`publiccode.yml`** files from the
organizations in the whitelist.

It then creates YAML files used by the
[Jekyll build chain](https://github.com/italia/developers.italia.it)
to generate the static pages of [developers.italia.it](https://developers.italia.it/).

[Elasticsearch 6.8](https://www.elastic.co/products/elasticsearch) is used to store
the data which be active and ready to accept connections before the crawler is started.

## Setup and deployment processes

The crawler can either run manually on the target machine or it can be deployed
in form of Docker container with
[its helm-chart](https://github.com/teamdigitale/devita-infra-kubernetes) in Kubernetes.

### Manually configure and build the crawler

1. `cd crawler`

2. Save the auth tokens to `domains.yml`.

3. Rename `config.toml.example` to `config.toml` and set the variables

   > **NOTE**: The application also supports environment variables in substitution
   > to config.toml file. Remember: "environment variables get higher priority than
   > the ones in configuration file"

4. Build the crawler binary with `make`

### Docker

The repository has a `Dockerfile`, used to build the production image,
and a `docker-compose.yml` file to facilitate the local deployment.

Before proceeding with the build, copy [`.env.example`](.env.example)
into `.env` and edit the environment variables as needed.

To build the crawler container run:

```shell
docker-compose up [-d] [--build]
```

where:

* *-d* execute the containers in background

* *--build* forces the containers build

To destroy the container, use:

```shell
docker-compose down
```

## Run the crawler

* Crawl mode (all item in whitelists): `bin/crawler crawl whitelist/*.yml`
  * `crawl` supports blacklists (see below for details). The crawler will try to
    match each repository URL in its list with the ones listed in blacklists and,
    if it does, it will print a warn log and skip all operation on it.
    Furthermore it will immediately remove the blacklisted repository from ES if
    it is present.

* One mode (single repository url): `bin/crawler one [repo url] whitelist/*.yml`
  * In this mode one single repository at the time will be evaluated. If the
    organization is present, its IPA code will be matched with the ones in
    whitelist otherwise it will be set to null and the `slug` will have a random
    code in the end (instead of the IPA code). Furthermore, the IPA code
    validation, which is a simple check within whitelists (to ensure that code
    belongs to the selected PA), will be skipped.
  * `one` supports blacklists (see below for details), whether `[repo url]` is
    present in one of the indicated blacklists, the crawler will exit immediately.
    Basically ignore all repository defined in list preventing the unauthorized
    loading in catalog.

* `bin/crawler updateipa` downloads IPA data and writes them into Elasticsearch

* `bin/crawler delete [URL]` deletes software from Elasticsearch using its code
   hosting URL specified in `publiccode.url`

* `bin/crawler download-whitelist` downloads organizations and repositories from
  the [onboarding portal repository](https://github.com/italia/developers-italia-onboarding)
  and saves them to a whitelist file

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

* [developers-italia-onboarding](https://github.com/italia/developers-italia-onboarding):
  the onboarding portal

## Authors

[Developers Italia](https://developers.italia.it) is a project by
[AgID](https://www.agid.gov.it/) and the
[Italian Digital Team](https://teamdigitale.governo.it/), which developed the
crawler and maintains this repository.
