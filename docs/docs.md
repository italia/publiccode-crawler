## developers-italia-backend

Backend & crawler for the OSS catalog of Developers Italia.

Table of contents:

- README
- Deploy Architecture (Docker & Containers)
- Files and folders description
- Data mapping (elasticsearch)
- Configuration files
- Crawler
- Jekyll files generation

- References

## Run crawler in cron job

Execute every 12 hours `0 */12 * * * make crawler > crawler.log`
