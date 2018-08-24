## developers-italia-backend

Backend & crawler for the OSS catalog of Developers Italia.

Table of contents:

- [README](../README.md)
- [Deploy Architecture (Docker & Containers)](deploy.md)
- [Files and folders description](fileAndFolders.md)
- [Elasticsearch details and data mapping](elasticsearch.md)
- [Configuration](fileAndFolders.md)
- [Crawler flow and steps](crawler.md)
- [Jekyll files generation](jekyll.md)

- [References](references.md)

### Run crawler in cron job

Execute every 12 hours `0 */12 * * * make crawler > crawler.log`
