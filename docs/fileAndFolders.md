## File and folders description

**Generic file and folders :**

.circleci: CircleCI folder

.git: git folder

cmd, crawler, httpclient, ipa, jekyll, metrics, vendor, version: golang packages and crawler source.

CRAWLER_DATADIR: the tree for publiccode.yml files:

- CRAWLER_DATADIR/repos/<host>/<org>/<repo>/publiccodes_publiccode.yml
- CRAWLER_DATADIR/repos/<host>/<org>/<repo>/publiccodes_metadata_publiccode.yml
- CRAWLER_DATADIR/repos/<host>/<org>/<repo>/gitClone/<git_clone_files>

./docker: the tree for docker config files

./elasticsearch: the tree for elasticsearch mappings and example files

vitality-ranges.yml: ranges for vitality index (Ref: vitality-ranges.md)

whitelist/generic.yml and whitelist/pa.yml: whitelists for Public Administrations and generic organizations ([Reference](references.md))

**Configuration files**

- Docker and crawler building configurations: .env
- Crawler specific configurations: config.toml
- Crawler domains basic auth: domains.yml
- Crawler vitality-ranges: vitality-ranges.yml
- Crawler whitelist: whitelist/generic.yml and whitelist/pa.yml
